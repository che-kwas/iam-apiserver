package apiserver

import (
	"github.com/che-kwas/iam-kit/errcode"
	"github.com/che-kwas/iam-kit/httputil"
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
}

func notFound() func(c *gin.Context) {
	return func(c *gin.Context) {
		httputil.WriteResponse(c, errors.WithCode(errcode.ErrNotFound, "Not found."), nil)
	}
}

// // Validation make sure users have the right resource permission and operation.
// func Validation() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := isAdmin(c); err != nil {
// 			switch c.FullPath() {
// 			case "/v1/users":
// 				if c.Request.Method != http.MethodPost {
// 					httputil.WriteResponse(c, errors.WithCode(errcode.ErrPermissionDenied, ""), nil)
// 					c.Abort()

// 					return
// 				}
// 			case "/v1/users/:name", "/v1/users/:name/change_password":
// 				username := c.GetString("username")
// 				if c.Request.Method == http.MethodDelete ||
// 					(c.Request.Method != http.MethodDelete && username != c.Param("name")) {
// 					httputil.WriteResponse(c, errors.WithCode(errcode.ErrPermissionDenied, ""), nil)
// 					c.Abort()

// 					return
// 				}
// 			default:
// 			}
// 		}

// 		c.Next()
// 	}
// }

// // isAdmin make sure the user is administrator.
// func isAdmin(c *gin.Context) error {
// 	username := c.GetString(middleware.UsernameKey)
// 	user, err := store.Client().Users().Get(c, username, metav1.GetOptions{})
// 	if err != nil {
// 		return errors.WithCode(errcode.ErrDatabase, err.Error())
// 	}

// 	if !user.IsAdmin {
// 		return errors.WithCode(errcode.ErrPermissionDenied, "user %s is not a administrator", username)
// 	}

// 	return nil
// }
