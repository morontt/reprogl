package http

import (
	base "net/http"
)

type LogResponseWriter interface {
	base.ResponseWriter
	Status() int
}

type Response struct {
	base.ResponseWriter
	StatusCode int
}

func (lrw *Response) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *Response) Write(body []byte) (int, error) {
	if _, ok := lrw.ResponseWriter.Header()["Cache-Control"]; !ok {
		lrw.ResponseWriter.Header().Set("Cache-Control", "private, no-cache")
	}

	return lrw.ResponseWriter.Write(body)
}

func (lrw *Response) Status() int {
	if lrw.StatusCode == 0 {
		return base.StatusOK
	}

	return lrw.StatusCode
}
