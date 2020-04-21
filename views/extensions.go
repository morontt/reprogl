package views

import "html/template"

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}
