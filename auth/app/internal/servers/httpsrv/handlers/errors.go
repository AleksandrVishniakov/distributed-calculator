package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request) (statusCode int, err error)

func Errors(next ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := next(w, r)
		if err != nil {
			var responseError *servers.ResponseError
			switch {
			case errors.As(err, &responseError):
				servers.WriteError(w, responseError.Code, responseError.Message)
				return
			default:
				servers.WriteError(w, statusCode, err.Error())
				return
			}
		}
	})
}
