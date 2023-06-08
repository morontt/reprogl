package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func pageOrRedirect(params map[string]string) (int, bool) {
	var page int
	pageString := params["page"]

	if pageString == "1" {
		return 1, true
	} else if pageString == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageString)
	}

	return page, false
}

func doESI(w http.ResponseWriter) {
	w.Header().Set("Surrogate-Control", "content=\"ESI/1.0\"")
}

func cacheControl(w http.ResponseWriter, age int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", age))
}

func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	jsonResult, _ := json.Marshal(data)
	w.Write(jsonResult)
}
