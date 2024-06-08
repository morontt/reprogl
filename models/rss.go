package models

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Rss struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RssChannel `xml:"channel"`
	Version string     `xml:"version,attr"`
}

type RssChannel struct {
	XMLName     xml.Name   `xml:"channel"`
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Link        string     `xml:"link"`
	Updated     RssTime    `xml:"pubDate"`
	Language    string     `xml:"language"`
	Author      string     `xml:"managingEditor"`
	Generator   string     `xml:"generator"`
	Docs        string     `xml:"docs"`
	Entry       []RssEntry `xml:"item"`
}

type RssEntry struct {
	ID          RssID          `xml:"guid"`
	Title       string         `xml:"title"`
	URL         string         `xml:"link"`
	Updated     RssTime        `xml:"pubDate"`
	Description RssDescription `xml:"description"`
	Enclosure   *RssEnclosure  `xml:"enclosure,omitempty"`
}

type RssID struct {
	XMLName   xml.Name `xml:"guid"`
	PermaLink string   `xml:"isPermaLink,attr"`
	ID        string   `xml:",innerxml"`
}

type RssDescription struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type RssEnclosure struct {
	XMLName  xml.Name `xml:"enclosure"`
	Url      string   `xml:"url,attr"`
	Length   int      `xml:"length,attr"`
	MimeType string   `xml:"type,attr"`
}

func (a *Rss) setChannelData(data FeedChannelData) {
	a.Version = "2.0"
	a.Channel = RssChannel{
		Title:       data.Title,
		Description: data.Description,
		Link:        data.Link,
		Language:    data.Language,
		Author:      fmt.Sprintf("%s (%s)", data.Email, data.Author),
		Generator:   data.Generator,
		Docs:        "http://www.rssboard.org/rss-specification",
	}

	for _, entry := range data.FeedItems {
		a.addFeedItem(entry)
	}
}

func (a *Rss) ContentType() string {
	return "application/atom+xml; charset=utf-8"
}

func (a *Rss) addFeedItem(entry *FeedItem) {
	a.Channel.Entry = append(a.Channel.Entry, RssEntry{
		ID: RssID{
			ID:        entry.GIUD(),
			PermaLink: "false",
		},
		Title:       entry.Title,
		URL:         entry.URL,
		Updated:     RssTime(entry.CreatedAt),
		Description: RssDescription{Text: entry.Text},
		Enclosure:   entry.GetRssEnclosure(),
	})

	if entry.CreatedAt.After(time.Time(a.Channel.Updated)) {
		a.Channel.Updated = RssTime(entry.CreatedAt)
	}
}
