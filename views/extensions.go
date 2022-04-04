package views

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"strings"
	"xelbot.com/reprogl/models"
)

var router *mux.Router

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

func absUrlGenerator(routeName string, pairs ...string) string {
	return "https://" + cfg.Host + urlGenerator(routeName, pairs...)
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
	return cfg.CDNBaseURL
}

func nl2br(s string) string {
	return strings.Replace(s, "\n", "<br/>", -1)
}

func authorName() string {
	return cfg.Author
}

func authorBio() string {
	return cfg.AuthorBio
}

func authorGithub() string {
	return fmt.Sprintf("https://github.com/%s", cfg.AuthorGithub)
}

func authorInstagram() string {
	return fmt.Sprintf("https://www.instagram.com/%s/", cfg.AuthorInsta)
}

func renderESI(routeName string, pairs ...string) template.HTML {
	s := fmt.Sprintf(
		"<esi:include src=\"%s\" onerror=\"continue\"/>",
		urlGenerator(routeName, pairs...),
	)

	return template.HTML(s)
}

func subString(input string, length int) (str string) {
	symbols := []rune(input)

	if len(symbols) >= length {
		str = string(symbols[:length-3]) + "..."
	} else {
		str = input
	}

	return
}
