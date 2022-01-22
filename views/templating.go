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
			"./templates/layout/base.gohtml",
		},
		"article.gohtml": {
			"./templates/article.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/layout/base.gohtml",
		},
		"index.gohtml": {
			"./templates/index.gohtml",
			"./templates/partials/menu.gohtml",
			"./templates/partials/sticky-header.gohtml",
			"./templates/partials/big-header.gohtml",
			"./templates/partials/footer.gohtml",
			"./templates/layout/base.gohtml",
		},
		"categories.gohtml": {
			"./templates/fragments/categories.gohtml",
		},
	}

	customFunctions := template.FuncMap{
		"raw":          rawHTML,
		"path":         urlGenerator,
		"tags":         tags,
		"topicPreview": topicPreview,
		"cdn":          cdnBase,
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

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	if container.IsDevMode() {
		err := LoadViewSet()
		if err != nil {
			return err
		}
	}

	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}

	var buf strings.Builder
	buf.Grow(defaultPageSize)

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(buf.String()))
	if err != nil {
		return err
	}

	return nil
}
