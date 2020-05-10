package views

import (
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/models"
)

var headerText string
var host string

func init() {
	cfg := config.Get()
	headerText = cfg.HeaderText
	host = cfg.Host
}

type Meta struct {
	Host       string
	HeaderText string
	titleParts []string
}

type HeaderLineInfo interface {
	HeaderLineDescription() string
	HeaderLineText() string
}

type ArticlePageData struct {
	Meta
	Article *models.Article
}

type IndexPageData struct {
	Meta
	HeaderInfo HeaderLineInfo
	PageNumber int
	Articles   models.ArticleList
}

func defaultMeta() Meta {
	return Meta{Host: host, HeaderText: headerText}
}

func (m *Meta) AppendTitle(str string) {
	m.titleParts = append(m.titleParts, str)
}

func (m *Meta) BrowserTitle() string {
	var title string
	if len(m.titleParts) > 0 {
		for _, p := range m.titleParts {
			title += p + " » "
		}
	}
	title += m.Host

	return title
}

func NewArticlePageData(article *models.Article) *ArticlePageData {
	return &ArticlePageData{Article: article, Meta: defaultMeta()}
}

func NewIndexPageData(articles models.ArticleList, page int) *IndexPageData {
	return &IndexPageData{Articles: articles, PageNumber: page, Meta: defaultMeta()}
}

func NewCategoryPageData(articles models.ArticleList, headerInfo HeaderLineInfo, page int) *IndexPageData {
	return &IndexPageData{Articles: articles, PageNumber: page, HeaderInfo: headerInfo, Meta: defaultMeta()}
}

func NewInfoPageData() *Meta {
	meta := defaultMeta()

	return &meta
}
