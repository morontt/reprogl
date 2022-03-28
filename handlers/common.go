package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xelbot.com/reprogl/container"
)

func pageOrRedirect(params map[string]string) (int, bool) {
	var page int
	pageString := params["page"]

	if pageString == "1" {
		return 1, true
	} else if pageString == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageString)
	}

	return page, false
}

func doESI(w http.ResponseWriter) {
	w.Header().Set("Surrogate-Control", "content=\"ESI/1.0\"")
}

func cacheControl(w http.ResponseWriter, age int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", age))
}

func urlBySlugGenerator(router *mux.Router) func(string) string {
	cfg := container.GetConfig()
	return func(slug string) string {
		url, _ := router.Get("article").URL("slug", slug)

		return "https://" + cfg.Host + url.String()
	}
}
