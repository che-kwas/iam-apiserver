// Package user implements the user handler.
package user

import (
	"log"

	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/meta"
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
	log.Printf("user.Create: %+v", user)

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
	params := &v1.User{}

	if err := c.ShouldBindJSON(params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	err := u.srv.Users().Update(c, c.Param("name"), params)
	httputil.WriteResponse(c, err, nil)
}

// ChangePasswordRequest defines the ChangePasswordRequest data format.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ChangePassword change the user's password by the user identifier.
func (u *UserController) ChangePassword(c *gin.Context) {
	var params ChangePasswordRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	err := u.srv.Users().ChangePassword(c, c.Param("name"), params.OldPassword, params.NewPassword)
	httputil.WriteResponse(c, err, nil)
}

// List lists the users in the storage.
func (u *UserController) List(c *gin.Context) {
	var opts meta.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	users, err := u.srv.Users().List(c, opts)
	httputil.WriteResponse(c, err, users)
}

// Delete deletes a user by the user identifier.
func (u *UserController) Delete(c *gin.Context) {
	err := u.srv.Users().Delete(c, c.Param("name"))
	httputil.WriteResponse(c, err, nil)
}

// DeleteCollection batch delete users by usernames.
func (u *UserController) DeleteCollection(c *gin.Context) {
	usernames := c.QueryArray("name")
	err := u.srv.Users().DeleteCollection(c, usernames)
	httputil.WriteResponse(c, err, nil)
}
