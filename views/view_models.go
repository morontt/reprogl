package views

import (
	"fmt"
	"reflect"
	"strings"
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
	Ogp          OpenGraph
	Canonical    string
}

type MetaName struct {
	Name    string
	Content string
}

type OpenGraph map[string]string

type HeaderLineInfo interface {
	HeaderLineDescription() string
	HeaderLineText() string
}

type DataWithFlashMessage interface {
	HasSuccessFlash() bool
	FlashSuccess() string
	SetSuccessFlash(string)
}

type flashObjectPart struct {
	flashSuccess string
}

type identityPart struct {
	identity security.Identity
}

type ArticlePageData struct {
	Meta
	Article      *models.Article
	CommentKey   string
	AcceptHeader string
	AuthorAvatar string
	flashObjectPart
	identityPart
}

type IndexPageData struct {
	Meta
	HeaderInfo HeaderLineInfo
	Paginator  *models.ArticlesPaginator
	flashObjectPart
}

type InfoPageData struct {
	Meta
	HeaderInfo HeaderLineInfo
	flashObjectPart
	Jobs container.JobHistory
}

type SidebarDummyArticle struct {
	IsArticle     bool
	RecentPostsID string
}

type StatisticsPageData struct {
	Meta
	Now          time.Time
	Commentators *models.CommentatorList
	DummyArticle SidebarDummyArticle

	MonthArticles   []models.ArticleStatItem
	AllTimeArticles []models.ArticleStatItem
}

type LoginPageData struct {
	Meta
	CsrfToken    string
	ErrorMessage string
	HasError     bool
}

type UnsubscribePageData struct {
	Meta
	Settings *models.EmailSubscription
	Avatar   string
	Success  bool
}

type AuthNavigation struct {
	identityPart
}

type MenuAuthData struct {
	User *models.User
	identityPart
}

type ProfilePageData struct {
	Meta
	User                  *models.User
	SubscriptionsSettings *models.EmailSubscription
	DummyArticle          SidebarDummyArticle
}

type OauthPendingPageData struct {
	Meta
	RequestId string
}

type FragmentCategoriesData struct {
	Categories *models.CategoryList
}

type FragmentCommentsData struct {
	Comments        models.CommentList
	EnabledComments bool
	identityPart
}

type FragmentRecentPostsData struct {
	RecentPosts *models.RecentPostList
}

func defaultMeta() Meta {
	cfg := container.GetConfig()

	ogp := make(OpenGraph)
	// og:title and og:url required for every page
	ogp["og:type"] = "website"
	ogp["og:image"] = cfg.CDNBaseURL + "/images/kravchik.jpg"
	ogp["og:image:width"] = "752"
	ogp["og:image:height"] = "376"
	ogp["og:locale"] = "ru_RU"

	return Meta{Host: cfg.Host, HeaderText: cfg.HeaderText, Ogp: ogp}
}

func (m *Meta) AppendTitle(str string) {
	m.titleParts = append(m.titleParts, str)
	m.SetOpenGraphProperty("og:title", str)
}

func (m *Meta) AppendName(name, content string) {
	m.MetaParts = append(m.MetaParts, MetaName{Name: name, Content: content})
	if name == "description" {
		m.SetOpenGraphProperty("og:description", content)
	}
}

func (m *Meta) SetCanonical(link string) {
	m.Canonical = link
	m.SetOpenGraphProperty("og:url", link)
}

func (m *Meta) SetOpenGraphProperty(property, content string) {
	if m.Ogp != nil {
		m.Ogp[property] = content
	} else {
		ogp := make(OpenGraph)
		ogp[property] = content
		m.Ogp = ogp
	}
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

func (fo *flashObjectPart) HasSuccessFlash() bool {
	return len(fo.flashSuccess) > 0
}

func (fo *flashObjectPart) FlashSuccess() string {
	return fo.flashSuccess
}

func (fo *flashObjectPart) SetSuccessFlash(msg string) {
	fo.flashSuccess = msg
}

func (ip *identityPart) SetIdentity(identity security.Identity) {
	ip.identity = identity
}

func (ip *identityPart) HasIdentity() bool {
	return !ip.identity.IsZero()
}

func (ip *identityPart) IsAdmin() bool {
	return ip.identity.IsAdmin()
}

func NewArticlePageData(article *models.Article, commentKey, accept string) *ArticlePageData {
	meta := defaultMeta()
	if article.Description.Valid {
		meta.AppendName("description", article.Description.String)
	}

	return &ArticlePageData{
		Article:      article,
		Meta:         meta,
		CommentKey:   commentKey,
		AcceptHeader: accept,
	}
}

func NewIndexPageData(paginator *models.ArticlesPaginator) *IndexPageData {
	meta := defaultMeta()
	meta.IsIndexPage = true

	return &IndexPageData{
		Paginator: paginator,
		Meta:      meta,
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
	meta.AppendTitle("Обо мне")

	return &InfoPageData{
		Meta: meta,
		Jobs: container.GetConfig().Jobs.Sort(),
	}
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
	meta.AppendName("description", "Страница логина. Описание тут не особо-то и нужно, но Yandex Webmaster настаивает")

	return &LoginPageData{
		Meta:         meta,
		CsrfToken:    token,
		ErrorMessage: errorMessage,
		HasError:     hasError,
	}
}

func NewUnsubscribePageData(settings *models.EmailSubscription, avatarLink string, success bool) *UnsubscribePageData {
	meta := defaultMeta()
	meta.AppendTitle("Отписка от email-уведомлений")

	return &UnsubscribePageData{
		Meta:     meta,
		Settings: settings,
		Avatar:   avatarLink,
		Success:  success,
	}
}

func NewAuthNavigationData() *AuthNavigation {
	return &AuthNavigation{}
}

func NewMenuAuthData(user *models.User) *MenuAuthData {
	return &MenuAuthData{User: user}
}

func NewProfilePageData(user *models.User, settings *models.EmailSubscription) *ProfilePageData {
	meta := defaultMeta()
	meta.AppendTitle("Профиль")

	return &ProfilePageData{
		Meta:         meta,
		User:         user,
		DummyArticle: SidebarDummyArticle{RecentPostsID: "0"},

		SubscriptionsSettings: settings,
	}
}

func NewOauthPendingPageData(requsetId string) *OauthPendingPageData {
	meta := defaultMeta()
	meta.AppendTitle("Ожидайте...")

	return &OauthPendingPageData{
		Meta:      meta,
		RequestId: requsetId,
	}
}

func (apd *ArticlePageData) AcceptWebp() bool {
	return strings.Contains(apd.AcceptHeader, "image/webp")
}

func (apd *ArticlePageData) AcceptAvif() bool {
	return strings.Contains(apd.AcceptHeader, "image/avif")
}
