package dto

import (
	"github.com/gin-gonic/gin"
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

func (e *ResponseError) Abort(c *gin.Context) {
	c.AbortWithStatusJSON(e.Code, e)
}