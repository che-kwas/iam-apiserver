package apiserver

import (
	"fmt"
	"iam-apiserver/internal/apiserver/store"

	"github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/che-kwas/iam-kit/middleware/auth"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func initRouter(g *gin.Engine) {
	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), notFound())

	jwtStrategy := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userController := user.NewUserController()
			userv1.POST("", userController.Create)

			userv1.Use(auto.AuthFunc())
			userv1.GET(":name", userController.Get)
			userv1.PUT(":name", userController.Update)
			userv1.PUT(":name/change-password", userController.ChangePassword)

			userv1.Use(isAdmin())
			userv1.GET("", userController.List)
			userv1.DELETE(":name", userController.Delete)
			userv1.DELETE("", userController.DeleteCollection)
		}
	}
}

func notFound() func(c *gin.Context) {
	return func(c *gin.Context) {
		httputil.WriteResponse(c, errors.WithCode(code.ErrNotFound, "Not found."), nil)
	}
}

func isAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString(middleware.UsernameKey)
		user, err := store.Client().Users().Get(c, username, metav1.GetOptions{})
		if err == nil && user.IsAdmin {
			c.Next()
			return
		}

		var msg string
		if err != nil {
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("user %s is not a administrator", username)
		}

		httputil.WriteResponse(c, errors.WithCode(code.ErrPermissionDenied, msg), nil)
		c.Abort()
		return
	}
}
