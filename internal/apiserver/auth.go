package apiserver

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/che-kwas/iam-kit/middleware/auth"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const JWTIssuer = "iam-apiserver"

var (
	ErrInvalidHeader = errors.New("invalid header")
)

type loginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		// fetch user from database
		user, err := store.Client().Users().Get(context.TODO(), username, metav1.GetOptions{})
		if err != nil {
			return false
		}

		// Compare the login password with the user password.
		if err := user.ComparePassword(password); err != nil {
			return false
		}

		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(context.TODO(), user, metav1.UpdateOptions{})

		return true
	})
}

func newJWTAuth() middleware.AuthStrategy {
	ginjwt, _ := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            "IAM",
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"),
		Authenticator:    authenticator(),
		LoginResponse:    loginResponse(),
		LogoutResponse:   logoutResponse(),
		RefreshResponse:  refreshResponse(),
		PayloadFunc:      payloadFunc(),
		IdentityHandler:  identityHandler(),
		IdentityKey:      middleware.UsernameKey,
		Authorizator:     authorizator(),
		Unauthorized:     unauthorizedHandler(),
		TokenLookup:      "header: Authorization",
		TokenHeadName:    "Bearer",
		SendCookie:       true,
		TimeFunc:         time.Now,
	})

	return auth.NewJWTStrategy(*ginjwt)
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth(), newJWTAuth())
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		login, err := parseWithHeader(c)
		if err != nil {
			return nil, ginjwt.ErrInvalidAuthHeader
		}

		user, err := store.Client().Users().Get(c, login.Username, metav1.GetOptions{})
		if err != nil {
			log.Error("Authentication failed: username error.")
			return nil, ginjwt.ErrFailedAuthentication
		}

		if err := user.ComparePassword(login.Password); err != nil {
			log.Error("Authentication failed: password error.")
			return nil, ginjwt.ErrFailedAuthentication
		}

		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(c, user, metav1.UpdateOptions{})

		return user, nil
	}
}

func parseWithHeader(c *gin.Context) (*loginInfo, error) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return nil, ErrInvalidHeader
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return nil, ErrInvalidHeader
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, ErrInvalidHeader
	}

	return &loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func logoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		c.JSON(http.StatusOK, nil)
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return loginResponse()
}

func payloadFunc() func(data interface{}) ginjwt.MapClaims {
	return func(data interface{}) ginjwt.MapClaims {
		claims := ginjwt.MapClaims{"iss": JWTIssuer}
		if u, ok := data.(*v1.User); ok {
			claims[ginjwt.IdentityKey] = u.Username
			claims["sub"] = u.Username
		}

		return claims
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := ginjwt.ExtractClaims(c)

		return claims[ginjwt.IdentityKey]
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.Infof("user `%s` is authenticated.", v)

			return true
		}

		return false
	}
}

func unauthorizedHandler() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"message": message,
		})
	}
}