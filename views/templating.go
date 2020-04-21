package views

import (
	"fmt"
	"html/template"
	"net/http"
	"xelbot.com/reprogl/config"
)

var templates map[string]*template.Template

var customFunctions = template.FuncMap{
	"raw": rawHTML,
}

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
	if config.IsDevMode() {
		err := LoadViewSet()
		if err != nil {
			return err
		}
	}

	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return tmpl.Execute(w, data)
}
