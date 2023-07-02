package views

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"xelbot.com/reprogl/container"
)

const defaultPageSize = 64 * 1024

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
}

func LoadViewSet() error {
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
	}

	customFunctions := template.FuncMap{
		"raw":             rawHTML,
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
		"flag_cnt":        flagCounterImage(true),
		"flag_cnt_mini":   flagCounterImage(false),
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

func WriteTemplate(w http.ResponseWriter, name string, data interface{}) error {
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
