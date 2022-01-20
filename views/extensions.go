package views

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"strings"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/models"
)

var router *mux.Router
var cdnBaseURL string

func init() {
	cfg := config.Get()
	cdnBaseURL = cfg.CDNBaseURL
}

func SetRouter(r *mux.Router) {
	router = r
}

func rawHTML(s string) template.HTML {
	return template.HTML(s)
}

func urlGenerator(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		panic(err)
	}

	return url.String()
}

func tags(tl models.TagList) template.HTML {
	var s string
	if len(tl) > 0 {
		s = "Теги: "
		links := make([]string, len(tl))
		for i, t := range tl {
			links[i] = fmt.Sprintf(
				"<a href=\"%s\">%s</a>",
				urlGenerator("tag-first", "slug", t.Slug),
				t.Name,
			)
		}

		s += strings.Join(links, ", ")
	}

	return template.HTML(s)
}

func topicPreview(s string) template.HTML {
	arr := strings.Split(s, "<!-- cut -->")
	if len(arr) > 1 {
		s = arr[0]
	}

	return template.HTML(s)
}

func cdnBase() string {
	return cdnBaseURL
}
