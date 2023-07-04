package views

import (
	"fmt"
	"reflect"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
)

type Meta struct {
	Host         string
	HeaderText   string
	MetaParts    []MetaName
	IsIndexPage  bool
	IsAuthorPage bool
	titleParts   []string
}

type MetaName struct {
	Name    string
	Content string
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

type SidebarDummyArticle struct {
	HasImage      bool
	IsArticle     bool
	RecentPostsID string
}

type StatisticsPageData struct {
	Meta
	Now          time.Time
	Commentators *models.CommentatorList
	DummyArticle SidebarDummyArticle
}

type LoginPageData struct {
	Meta
	CsrfToken string
}

type AuthNavigation struct {
	Authenticated bool
}

type FragmentCategoriesData struct {
	Categories *models.CategoryList
}

type FragmentCommentsData struct {
	Comments        *models.CommentList
	EnabledComments bool
}

type FragmentRecentPostsData struct {
	RecentPosts *models.RecentPostList
}

func defaultMeta() Meta {
	cfg := container.GetConfig()

	return Meta{Host: cfg.Host, HeaderText: cfg.HeaderText}
}

func (m *Meta) AppendTitle(str string) {
	m.titleParts = append(m.titleParts, str)
}

func (m *Meta) AppendName(name, content string) {
	m.MetaParts = append(m.MetaParts, MetaName{Name: name, Content: content})
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
	meta := defaultMeta()
	if article.Description.Valid {
		meta.AppendName("description", article.Description.String)
	}

	return &ArticlePageData{Article: article, Meta: meta, CommentKey: commentKey}
}

func NewIndexPageData(paginator *models.ArticlesPaginator) *IndexPageData {
	meta := defaultMeta()
	meta.IsIndexPage = true

	return &IndexPageData{Paginator: paginator, Meta: meta}
}

func NewCategoryPageData(paginator *models.ArticlesPaginator, headerInfo HeaderLineInfo) *IndexPageData {
	var browserTitle, metaDescription string
	meta := defaultMeta()
	meta.IsIndexPage = true

	switch reflect.TypeOf(headerInfo).String() {
	case "*models.Category":
		browserTitle = fmt.Sprintf("Категория \"%s\"", headerInfo.HeaderLineText())
		metaDescription = fmt.Sprintf("Записи из категории \"%s\"", headerInfo.HeaderLineText())
	case "*models.Tag":
		browserTitle = fmt.Sprintf("Тег \"%s\"", headerInfo.HeaderLineText())
		metaDescription = fmt.Sprintf("Записи по тегу \"%s\"", headerInfo.HeaderLineText())
	}

	if paginator.CurrentPage > 1 {
		browserTitle += fmt.Sprintf(". Страница %d", paginator.CurrentPage)
		metaDescription += fmt.Sprintf(". Страница %d", paginator.CurrentPage)
	}
	meta.AppendTitle(browserTitle)
	meta.AppendName("description", metaDescription)
	meta.AppendName("robots", "noindex, follow")

	return &IndexPageData{Paginator: paginator, HeaderInfo: headerInfo, Meta: meta}
}

func NewInfoPageData() *InfoPageData {
	meta := defaultMeta()
	meta.IsAuthorPage = true
	meta.AppendName("description", "Персональный блог Харченко Александра. Общая информация.")
	meta.AppendTitle("Информация")

	return &InfoPageData{Meta: meta}
}

func NewStatisticsPageData() *StatisticsPageData {
	meta := defaultMeta()
	meta.AppendName("description", "Статистика посещений и комментариев.")
	meta.AppendTitle("Статистика")

	return &StatisticsPageData{
		Meta:         meta,
		Now:          time.Now(),
		DummyArticle: SidebarDummyArticle{RecentPostsID: "0"},
	}
}

func NewLoginPageData(token string) *LoginPageData {
	meta := defaultMeta()
	meta.AppendTitle("Вход")

	return &LoginPageData{
		Meta:      meta,
		CsrfToken: token,
	}
}

func NewAuthNavigationData(has bool) *AuthNavigation {
	return &AuthNavigation{
		Authenticated: has,
	}
}
