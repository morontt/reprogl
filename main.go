package main

import (
	"log"
	"net/http"
	"reprogl/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
