package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenToken(kid string, secret string, expire time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expire.Unix(),
		"iat": time.Now().Unix(),
	})

	token.Header["kid"] = kid
	return token.SignedString([]byte(secret))
}
