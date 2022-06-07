package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenToken(kid string, secret string, expire time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"kid": kid,
		"exp": expire.Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}
