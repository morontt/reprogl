package views

import (
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

var headerText string
var host string

func init() {
	cfg := container.GetConfig()
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
	Paginator  *models.ArticlesPaginator
}

type FragmentCategoriesData struct {
	Categories *models.CategoryList
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

func NewIndexPageData(paginator *models.ArticlesPaginator) *IndexPageData {
	return &IndexPageData{Paginator: paginator, Meta: defaultMeta()}
}

func NewCategoryPageData(paginator *models.ArticlesPaginator, headerInfo HeaderLineInfo) *IndexPageData {
	return &IndexPageData{Paginator: paginator, HeaderInfo: headerInfo, Meta: defaultMeta()}
}

func NewInfoPageData() *Meta {
	meta := defaultMeta()

	return &meta
}
