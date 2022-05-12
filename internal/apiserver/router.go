package apiserver

import (
	"github.com/che-kwas/iam-kit/errcode"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func initRouter(g *gin.Engine) {
	jwtStrategy := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), notFound())
}

func notFound() func(c *gin.Context) {
	return func(c *gin.Context) {
		httputil.WriteResponse(c, errors.WithCode(errcode.ErrNotFound, "Not found."), nil)
	}
}
