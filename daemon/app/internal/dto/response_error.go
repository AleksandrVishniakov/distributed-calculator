package dto

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ResponseError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewResponseError(code int, message string) *ResponseError {
	return &ResponseError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

func (e *ResponseError) Abort(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println("fail to return error:", *e)
	}
}
