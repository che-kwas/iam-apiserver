// Package user implements the user handler.
package user

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/service"
)

// UserController handles requests for user resource.
type UserController struct {
	srv service.Service
}

// NewUserController creates a user handler.
func NewUserController() *UserController {
	return &UserController{
		srv: service.NewService(),
	}
}

// Create creates a new user.
func (u *UserController) Create(c *gin.Context) {
	user := &v1.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	err := u.srv.Users().Create(c, user)
	httputil.WriteResponse(c, err, user)
}

// Get gets the user by the user identifier.
func (u *UserController) Get(c *gin.Context) {
	user, err := u.srv.Users().Get(c, c.Param("name"))

	httputil.WriteResponse(c, err, user)
}

// Update updates the user info by the user identifier.
func (u *UserController) Update(c *gin.Context) {
	user := &v1.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	err := u.srv.Users().Update(c, c.Param("name"), user)
	httputil.WriteResponse(c, err, user)
}
