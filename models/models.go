package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

var RecordNotFound = errors.New("models: no matching record found")

const (
	DisabledComments = "d"
	EnabledComments  = "e"
)

type ArticleBasePart struct {
	ID           int
	Title        string
	Slug         string
	Text         string
	Preview      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CategoryName string
	CategorySlug string
	Hidden       bool

	CommentsCount int
}

type Article struct {
	ArticleBasePart
	Description sql.NullString
	FeaturedImage
	Tags TagList

	DisabledComments bool
	RecentPostsID    string
	Views            int
	LjItemID         sql.NullInt32
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
	LeftKey  sql.NullInt32
	RightKey sql.NullInt32
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

type ArticleForComment struct {
	ID     int
	Slug   string
	Hidden bool
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

func (a *ArticleBasePart) IsArticle() bool {
	return true
}

func (a *Article) DisabledCommentsFlag() (flag string) {
	if a.DisabledComments {
		flag = DisabledComments
	} else {
		flag = EnabledComments
	}

	return
}
