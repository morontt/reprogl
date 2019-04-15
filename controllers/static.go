package controllers

import (
	"fmt"
	"net/http"
)

func InfoAction(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Info page")
}
