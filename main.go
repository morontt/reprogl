package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"reprogl/config"
	"reprogl/controllers"
	"reprogl/middlewares"
)

func main() {
	handleError(config.Load())

	handler := middlewares.AccessLog(getRoutes())

	cfg := config.Get()

	listener, _ := net.Listen("tcp", ":"+cfg.Port)
	handleError(fcgi.Serve(listener, handler))
}

func getRoutes() http.Handler {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/{page:[0-9]*}", controllers.IndexAction).Name("blog-page")
	siteMux.HandleFunc("/category/{slug}/{page:[0-9]*}", controllers.CategoryAction).Name("category")
	siteMux.HandleFunc("/tag/{slug}/{page:[0-9]*}", controllers.TagAction).Name("tag")
	siteMux.HandleFunc("/info", controllers.InfoAction).Name("info_page")

	return siteMux
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
