package models

import (
	"database/sql"
	"errors"
	"time"
	"xelbot.com/reprogl/config"
)

var RecordNotFound = errors.New("models: no matching record found")

var cdnBaseURL string

func init() {
	cfg := config.Get()
	cdnBaseURL = cfg.CDNBaseURL
}

type FeaturedImage struct {
	ImagePath        sql.NullString
	ImageDescription sql.NullString
}

type Article struct {
	ID            int
	Title         string
	Slug          string
	Text          string
	Description   sql.NullString
	CreatedAt     time.Time
	CommentsCount int
	CategoryName  string
	CategorySlug  string
	FeaturedImage
}

type ArticleListItem struct {
	ID           int
	Title        string
	Slug         string
	Text         string
	CreatedAt    time.Time
	CategoryName string
	CategorySlug string
	FeaturedImage
}

func (i *FeaturedImage) HasImage() bool {
	return i.ImagePath.Valid
}

func (i *FeaturedImage) ImageURL() string {
	var url string
	if i.ImagePath.Valid {
		url = cdnBaseURL + "/uploads/" + i.ImagePath.String
	} else {
		url = ""
	}

	return url
}

type ArticleList []*ArticleListItem
