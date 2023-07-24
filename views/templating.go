package views

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/security"
	"xelbot.com/reprogl/session"
)

const defaultPageSize = 64 * 1024

var (
	templates          map[string]*template.Template
	templatesMapLocker sync.Mutex
)

func init() {
	templates = make(map[string]*template.Template)
}

func LoadViewSet() error {
	templatesMapLocker.Lock()
	defer templatesMapLocker.Unlock()

	templatesMap := map[string][]string{
		"info.gohtml": {
			"./templates/info.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/partials/sticky-header.gohtml",
			"./templates/partials/big-header.gohtml",
			"./templates/partials/footer.gohtml",
			"./templates/partials/social-icons.gohtml",
			"./templates/partials/info-static.gohtml",
			"./templates/layout/svg-sprites.gohtml",
			"./templates/layout/base.gohtml",
		},
		"article.gohtml": {
			"./templates/article.gohtml",
			"./templates/partials/author-info.gohtml",
			"./templates/partials/comment-form.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/partials/sidebar.gohtml",
			"./templates/partials/social-icons.gohtml",
			"./templates/layout/svg-sprites.gohtml",
			"./templates/layout/base.gohtml",
		},
		"statistics.gohtml": {
			"./templates/statistics.gohtml",
			"./templates/partials/author-info.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/partials/sidebar.gohtml",
			"./templates/partials/social-icons.gohtml",
			"./templates/layout/svg-sprites.gohtml",
			"./templates/layout/base.gohtml",
		},
		"index.gohtml": {
			"./templates/index.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/partials/sticky-header.gohtml",
			"./templates/partials/big-header.gohtml",
			"./templates/partials/footer.gohtml",
			"./templates/partials/social-icons.gohtml",
			"./templates/layout/svg-sprites.gohtml",
			"./templates/layout/base.gohtml",
		},
		"categories.gohtml": {
			"./templates/fragments/categories.gohtml",
		},
		"comments.gohtml": {
			"./templates/fragments/comments.gohtml",
		},
		"recent-posts.gohtml": {
			"./templates/fragments/recent-posts.gohtml",
		},
		"login.gohtml": {
			"./templates/auth/login.gohtml",
		},
		"auth-navigation.gohtml": {
			"./templates/fragments/auth-navigation.gohtml",
		},
	}

	customFunctions := template.FuncMap{
		"raw":             rawHTML,
		"is_dev":          isDev,
		"path":            urlGenerator,
		"abs_path":        absUrlGenerator,
		"render_esi":      renderESI,
		"tags":            tags,
		"cdn":             cdnBase,
		"nl2br":           nl2br,
		"author_name":     authorName,
		"author_bio":      authorBio,
		"author_github":   authorGithub,
		"author_telegram": authorTelegram,
		"substr":          subString,
		"time_tag":        timeTag,
		"asset":           assetTag,
		"go_version":      goVersion,
		"cnt_comments":    commentsCountString,
		"cnt_times":       timesCountString,
		"flag_cnt":        flagCounterImage(true),
		"flag_cnt_mini":   flagCounterImage(false),
		"emojiFlag":       emojiFlag,

		"articleStyles":    articleStyles,
		"statisticsStyles": statisticsStyles,
		"indexStyles":      indexStyles,
		"infoStyles":       infoStyles,
	}

	for key, files := range templatesMap {
		tmpl, err := template.New(key).Funcs(customFunctions).ParseFiles(files...)
		if err != nil {
			return err
		}

		templates[key] = tmpl
	}

	return nil
}

func RenderTemplate(name string, data interface{}) (string, error) {
	if container.IsDevMode() {
		err := LoadViewSet()
		if err != nil {
			return "", err
		}
	}

	tmpl, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("the template %s does not exist", name)
	}

	var buf strings.Builder
	buf.Grow(defaultPageSize)

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func WriteTemplate(w http.ResponseWriter, name string, data any) error {
	content, err := RenderTemplate(name, data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Surrogate-Control", "content=\"ESI/1.0\"")
	_, err = w.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func WriteTemplateWithContext(ctx context.Context, w http.ResponseWriter, name string, data any) error {
	if flashObjectPart, ok := data.(DataWithFlashMessage); ok {
		if flashSuccessMessage, found := session.Pop[string](ctx, session.FlashSuccessKey); found {
			flashObjectPart.SetSuccessFlash(flashSuccessMessage)
		}
	}

	if identityPart, ok := data.(security.IdentityAware); ok {
		identity, _ := session.GetIdentity(ctx)
		identityPart.SetIdentity(identity)
	}

	return WriteTemplate(w, name, data)
}
