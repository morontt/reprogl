package views

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

func LoadViewSet() error {
	templates = make(map[string]*template.Template)

	files := []string{
		"./templates/static/info.gohtml",
		"./templates/base.gohtml",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	templates["static/info"] = tmpl

	return nil
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return tmpl.Execute(w, data)
}
