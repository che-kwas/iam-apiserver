package apiserver

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/che-kwas/iam-kit/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

const (
	JWTIssuer = "iam-apiserver"

	ConfKeyJWT           = "jwt"
	DefaultJWTTimeout    = time.Duration(24 * time.Second)
	DefaultJWTMaxRefresh = time.Duration(24 * time.Second)
)

var (
	ErrInvalidHeader = errors.New("invalid header")
)

// JWTOptions defines options for building a GinJWTMiddleware.
type JWTOptions struct {
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration `mapstructure:"max-refresh"`
}
type loginInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		// fetch user from database
		user, err := store.Client().Users().Get(context.TODO(), username)
		if err != nil {
			return false
		}

		if !user.VerifyPassword(password) {
			return false
		}

		now := time.Now()
		user.LoginedAt = &now
		_ = store.Client().Users().Update(context.TODO(), user)

		return true
	})
}

func newJWTAuth() middleware.AuthStrategy {
	opts := getJWTOptions()

	ginjwt, _ := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            "IAM",
		SigningAlgorithm: "HS256",
		Key:              []byte(opts.Key),
		Timeout:          opts.Timeout,
		MaxRefresh:       opts.MaxRefresh,
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

		user, err := store.Client().Users().Get(c, login.Username)
		if err != nil {
			return nil, ginjwt.ErrFailedAuthentication
		}

		if !user.VerifyPassword(login.Password) {
			return nil, ginjwt.ErrFailedAuthentication
		}

		now := time.Now()
		user.LoginedAt = &now
		_ = store.Client().Users().Update(c, user)

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
		if _, ok := data.(string); ok {
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

func getJWTOptions() *JWTOptions {
	opts := &JWTOptions{
		Timeout:    DefaultJWTTimeout,
		MaxRefresh: DefaultJWTMaxRefresh,
	}

	_ = viper.UnmarshalKey(ConfKeyJWT, opts)
	return opts
}
