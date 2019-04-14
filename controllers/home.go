package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func IndexAction(w http.ResponseWriter, r *http.Request) {
	var page int
	vars := mux.Vars(r)
	pageString := vars["page"]

	if pageString == "1" {
		http.Redirect(w, r, "/", 301)

		return
	} else if pageString == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageString)
	}

	fmt.Fprintf(w, "Articles, page %d", page)
}
