package main

import (
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
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", controllers.Index)

	return siteMux
}
