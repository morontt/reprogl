package handlers

import (
	"net/http"
	"strconv"
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
