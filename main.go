package main

import (
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/http/fcgi"
	"reprogl/controllers"
	"reprogl/middlewares"
)

func main() {
	handler := middlewares.AccessLog(getRoutes())

	listener, _ := net.Listen("tcp", ":8080")
	fcgi.Serve(listener, handler)
}

func getRoutes() http.Handler {
	siteMux := mux.NewRouter()
	siteMux.HandleFunc("/{page:[0-9]*}", controllers.IndexAction)

	return siteMux
}
