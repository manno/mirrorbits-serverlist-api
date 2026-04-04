package main

import (
	"net/http"
)

func NewRouter(context *appContext) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", Logger(appHandler{context, Index}, "Index"))
	return mux
}
