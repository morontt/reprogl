package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
	"xelbot.com/reprogl/container"
)

var RecordNotFound = errors.New("models: no matching record found")

var cdnBaseURL string

func init() {
	cfg := container.GetConfig()
	cdnBaseURL = cfg.CDNBaseURL
}

type FeaturedImage struct {
	ImagePath        sql.NullString
	ImageDescription sql.NullString
}

type ArticleBasePart struct {
	ID           int
	Title        string
	Slug         string
	Text         string
	Preview      string
	CreatedAt    time.Time
	CategoryName string
	CategorySlug string
}

type Article struct {
	ArticleBasePart
	Description   sql.NullString
	CommentsCount int
	FeaturedImage
	Tags TagList
}

type ArticleListItem struct {
	ArticleBasePart
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

type CategoryList []*Category

type Tag struct {
	ID   int
	Name string
	Slug string
}

type TagList []*Tag

type RecentPost struct {
	Title string
	Slug  string
}

type RecentPostList []*RecentPost

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

func (cat *Category) NameWithTreePrefix() string {
	var prefix string
	if cat.Depth > 1 {
		prefix = strings.Repeat("..", cat.Depth-1)
	}

	return prefix + cat.Name
}

func (tag *Tag) HeaderLineDescription() string {
	return "Записи с отметкой"
}

func (tag *Tag) HeaderLineText() string {
	return tag.Name
}

func (a *ArticleBasePart) HasPreview() bool {
	return strings.Contains(a.Text, "<!-- cut -->")
}

func (a *ArticleBasePart) IdString() string {
	return strconv.Itoa(a.ID)
}
