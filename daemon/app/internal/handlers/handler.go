package handlers

import "net/http"

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", nil)
	mux.HandleFunc("POST /subtract", nil)
	mux.HandleFunc("POST /multiply", nil)
	mux.HandleFunc("POST /divide", nil)

	return mux
}
