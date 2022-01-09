package models

import (
	"database/sql"
	"errors"
	"strings"
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
	Tags TagList
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
	Tags TagList
}

type ArticleList []*ArticleListItem

type Category struct {
	ID       int
	Name     string
	Slug     string
	LeftKey  sql.NullInt64
	RightKey sql.NullInt64
	Depth    int
}

type Tag struct {
	ID   int
	Name string
	Slug string
}

type TagList []*Tag

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

func (cat *Category) HeaderLineDescription() string {
	return "Записи из категории"
}

func (cat *Category) HeaderLineText() string {
	return cat.Name
}

func (tag *Tag) HeaderLineDescription() string {
	return "Записи с отметкой"
}

func (tag *Tag) HeaderLineText() string {
	return tag.Name
}

func (a *ArticleListItem) HasPreview() bool {
	return strings.Contains(a.Text, "<!-- cut -->")
}
