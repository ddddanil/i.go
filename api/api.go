package api

import "net/http"

type Api struct{}

func NewApi() http.Handler {
	mux := http.NewServeMux()
	return mux
}
