package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	UserIdContextKey = "user_id"
)

func NewJWTAuthMiddleware(tokensManager *jwt.TokenGenerator) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			accessToken, err := GetAccessToken(c.Request)
			if err != nil {
				dto.NewResponseError(http.StatusUnauthorized, err.Error()).Abort(c)
				return
			}

			userID, _, err := tokensManager.ParseToken(accessToken)
			if err != nil {
				dto.NewResponseError(http.StatusUnauthorized, err.Error()).Abort(c)
				return
			}

			log.Println("userID:", userID)

			c.Set(UserIdContextKey, userID)

			c.Next()
		}
	}
}

func GetAccessToken(r *http.Request) (string, error) {
	headerParts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header length")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth method")
	}

	return headerParts[1], nil
}
