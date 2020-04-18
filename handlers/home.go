package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func IndexAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles, page %d", page)
}

func CategoryAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles by category, page %d", page)
}

func TagAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, needsRedirect := pageOrRedirect(vars)
	if needsRedirect {
		http.Redirect(w, r, "/", 301)

		return
	}

	fmt.Fprintf(w, "Articles by tag, page %d", page)
}
