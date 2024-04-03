package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenGenerator struct {
	Signature []byte
	TokenTTL  time.Duration
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID uint64 `json:"userID"`
	Login  string `json:"login"`
}

func (g *TokenGenerator) NewToken(
	userID uint64,
	login string,
) (token string, err error) {
	accessTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(g.TokenTTL).Unix(),
		},
		UserID: userID,
		Login:  login,
	}).SignedString(g.Signature)

	if err != nil {
		return "", err
	}

	return accessTokenStr, nil
}

func (g *TokenGenerator) ParseToken(token string) (userID uint64, login string, err error) {
	accessToken, err := jwt.ParseWithClaims(
		token, &tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return g.Signature, nil
		})

	if err != nil {
		return 0, "", err
	}

	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("undefined token claims type")
	}

	return claims.UserID, claims.Login, nil
}
