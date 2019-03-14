package main

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"reprogl/controllers"
	"reprogl/middlewares"
)

func main() {
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", controllers.Index)

	handler := middlewares.AccessLog(siteMux)

	listener, _ := net.Listen("tcp", ":8080")
	fcgi.Serve(listener, handler)
}
