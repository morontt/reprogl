package main

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"reprogl/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Index)

	listener, _ := net.Listen("tcp", ":8080")
	fcgi.Serve(listener, nil)
}
