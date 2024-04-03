package servers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/parser"
)

type ResponseError struct {
	Code          int           `json:"code"`
	Message       string        `json:"message"`
	Timestamp     time.Time     `json:"timestamp"`
	DeveloperCode DeveloperCode `json:"developerCode"`
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("reponse error: %d %s", r.Code, r.Message)
}

func NewResponseError(
	code int,
	message string,
) *ResponseError {
	var devCode DeveloperCode

	devCode, ok := DevCodeFromMsg(message)
	if !ok {
		devCode = DeveloperCode(code * 1000)
	}

	return &ResponseError{
		Code:          code,
		Message:       message,
		Timestamp:     time.Now(),
		DeveloperCode: devCode,
	}
}

func WriteError(
	w http.ResponseWriter,
	code int,
	message string,
) {
	err := parser.EncodeResponse(w, NewResponseError(code, message), code)
	if err != nil {
		slog.Error(
			"error parsing response error",
			slog.String("error", err.Error()),
		)
	}
}
