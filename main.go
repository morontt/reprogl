package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
	"xelbot.com/reprogl/config"
	"xelbot.com/reprogl/handlers"
	"xelbot.com/reprogl/middlewares"
	"xelbot.com/reprogl/views"
)

func main() {
	handleError(config.Load())

	handler := middlewares.AccessLog(getRoutes())

	cfg := config.Get()

	handleError(views.LoadViewSet())

	server := &http.Server{
		Handler:      handler,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	handleError(server.ListenAndServe())
}

func getRoutes() http.Handler {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/{page:[0-9]*}", handlers.IndexAction).Name("blog-page")
	siteMux.HandleFunc("/category/{slug}/{page:[0-9]*}", handlers.CategoryAction).Name("category")
	siteMux.HandleFunc("/tag/{slug}/{page:[0-9]*}", handlers.TagAction).Name("tag")
	siteMux.HandleFunc("/info", handlers.InfoAction).Name("info_page")
	siteMux.HandleFunc("/robots.txt", handlers.RobotsTXTAction)

	return siteMux
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
