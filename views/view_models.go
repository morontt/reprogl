package views

import (
	"fmt"
	"reflect"
	"time"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/models"
	"xelbot.com/reprogl/security"
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
	Article     *models.Article
	CommentKey  string
	HasIdentity bool
}

type IndexPageData struct {
	Meta
	HeaderInfo   HeaderLineInfo
	Paginator    *models.ArticlesPaginator
	FlashSuccess string
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
	CsrfToken    string
	ErrorMessage string
	HasError     bool
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
	HasIdentity     bool
	Identity        security.Identity
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

func (ipd *IndexPageData) HasSuccessFlash() bool {
	return len(ipd.FlashSuccess) > 0
}

func NewArticlePageData(article *models.Article, commentKey string, hasIdentity bool) *ArticlePageData {
	meta := defaultMeta()
	if article.Description.Valid {
		meta.AppendName("description", article.Description.String)
	}

	return &ArticlePageData{
		Article:     article,
		Meta:        meta,
		CommentKey:  commentKey,
		HasIdentity: hasIdentity,
	}
}

func NewIndexPageData(paginator *models.ArticlesPaginator, flashSuccess string) *IndexPageData {
	meta := defaultMeta()
	meta.IsIndexPage = true

	return &IndexPageData{
		Paginator:    paginator,
		Meta:         meta,
		FlashSuccess: flashSuccess,
	}
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

func NewLoginPageData(token, errorMessage string, hasError bool) *LoginPageData {
	meta := defaultMeta()
	meta.AppendTitle("Вход")

	return &LoginPageData{
		Meta:         meta,
		CsrfToken:    token,
		ErrorMessage: errorMessage,
		HasError:     hasError,
	}
}

func NewAuthNavigationData(has bool) *AuthNavigation {
	return &AuthNavigation{
		Authenticated: has,
	}
}
