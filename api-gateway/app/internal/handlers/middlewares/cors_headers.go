package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusOK)
			return
		} else {
			c.Next()
		}
	}
}
