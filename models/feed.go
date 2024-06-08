package models

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"time"
)

const (
	RssFeedType = iota
	AtomFeedType
)

type SitemapTime time.Time

type RssTime time.Time

type FeedInterface interface {
	setChannelData(FeedChannelData)
	ContentType() string
	addFeedItem(*FeedItem)
}

type FeedGeneratorInterface interface {
	ContentType() string
	AsXML() ([]byte, error)
}

type FeedItem struct {
	ID        int
	Title     string
	Slug      string
	URL       string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	FeaturedImage
}

type FeedItemList []*FeedItem

type FeedChannelData struct {
	Title       string
	Link        string
	Description string
	Language    string
	Charset     string
	Author      string
	Email       string
	Generator   string
	FeedItems   FeedItemList
}

type Feed[F FeedInterface] struct {
	value F
}

func (ct SitemapTime) MarshalText() ([]byte, error) {
	t := time.Time(ct)

	return []byte(t.Format(time.RFC3339)), nil
}

func (ct RssTime) MarshalText() ([]byte, error) {
	t := time.Time(ct)

	return []byte(t.Format(time.RFC1123Z)), nil
}

func (d FeedChannelData) GIUD() string {
	return byStringGIUD(d.Link)
}

func (i FeedItem) GIUD() string {
	return byStringGIUD(i.Slug)
}

func CreateFeed[T FeedInterface](t T, d FeedChannelData) *Feed[T] {
	feed := Feed[T]{value: t}
	feed.value.setChannelData(d)

	return &feed
}

func (f Feed[T]) ContentType() string {
	return f.value.ContentType()
}

func (f Feed[T]) AsXML() ([]byte, error) {
	return xml.MarshalIndent(f.value, "", "  ")
}

func byStringGIUD(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	sum := hash.Sum(nil)

	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}
