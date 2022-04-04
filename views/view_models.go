package views

import (
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

var cfg container.AppConfig

func init() {
	cfg = container.GetConfig()
}

type Meta struct {
	Host            string
	HeaderText      string
	MetaDescription string
	IsIndexPage     bool
	IsAuthorPage    bool
	titleParts      []string
}

type HeaderLineInfo interface {
	HeaderLineDescription() string
	HeaderLineText() string
}

type ArticlePageData struct {
	Meta
	Article    *models.Article
	CommentKey string
}

type IndexPageData struct {
	Meta
	HeaderInfo HeaderLineInfo
	Paginator  *models.ArticlesPaginator
}

type InfoPageData struct {
	Meta
	HeaderInfo HeaderLineInfo
}

type FragmentCategoriesData struct {
	Categories *models.CategoryList
}

type FragmentCommentsData struct {
	Comments *models.CommentList
}

type FragmentRecentPostsData struct {
	RecentPosts *models.RecentPostList
}

func defaultMeta() Meta {
	return Meta{Host: cfg.Host, HeaderText: cfg.HeaderText}
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

func NewArticlePageData(article *models.Article, commentKey string) *ArticlePageData {
	return &ArticlePageData{Article: article, Meta: defaultMeta(), CommentKey: commentKey}
}

func NewIndexPageData(paginator *models.ArticlesPaginator) *IndexPageData {
	meta := defaultMeta()
	meta.IsIndexPage = true

	return &IndexPageData{Paginator: paginator, Meta: meta}
}

func NewCategoryPageData(paginator *models.ArticlesPaginator, headerInfo HeaderLineInfo) *IndexPageData {
	meta := defaultMeta()
	meta.IsIndexPage = true

	return &IndexPageData{Paginator: paginator, HeaderInfo: headerInfo, Meta: meta}
}

func NewInfoPageData() *InfoPageData {
	meta := defaultMeta()
	meta.IsAuthorPage = true
	meta.MetaDescription = "Персональный блог Харченко Александра. Общая информация."

	return &InfoPageData{Meta: meta}
}
