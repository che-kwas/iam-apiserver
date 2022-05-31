// Package user is the user controller.
package user

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/service"
)

// UserController handles requests for user resource.
type UserController struct {
	srv service.Service
	log *logger.Logger
}

// NewUserController creates a user controller.
func NewUserController() *UserController {
	return &UserController{
		srv: service.NewService(),
		log: logger.L(),
	}
}

// Create creates a new user.
func (u *UserController) Create(c *gin.Context) {
	u.log.X(c).Info("user create")
	user := &v1.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	u.log.Debugf("user create params: %+v", user)

	err := u.srv.Users().Create(c, user)
	httputil.WriteResponse(c, err, user)
}

// Get gets the user by the user identifier.
func (u *UserController) Get(c *gin.Context) {
	u.log.X(c).Info("user get")
	username := c.Param("name")
	u.log.Debugf("user get params: %s", username)

	user, err := u.srv.Users().Get(c, username)

	httputil.WriteResponse(c, err, user)
}

// Update updates the user info by the user identifier.
func (u *UserController) Update(c *gin.Context) {
	u.log.X(c).Info("user update")
	params := &v1.User{}

	if err := c.ShouldBindJSON(params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	u.log.Debugf("user update params: %+v", params)

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
	u.log.X(c).Info("user change-password")
	var params ChangePasswordRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	u.log.Debugf("user change-password params: %+v", params)

	err := u.srv.Users().ChangePassword(c, c.Param("name"), params.OldPassword, params.NewPassword)
	httputil.WriteResponse(c, err, nil)
}

// List lists the users in the storage.
func (u *UserController) List(c *gin.Context) {
	u.log.X(c).Info("user list")
	var opts meta.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	u.log.Debugf("user list params: %+v", opts)

	users, err := u.srv.Users().List(c, opts)
	httputil.WriteResponse(c, err, users)
}

// Delete deletes a user by the user identifier.
func (u *UserController) Delete(c *gin.Context) {
	u.log.X(c).Info("user delete")
	username := c.Param("name")
	u.log.Debugf("user delete params: %s", username)

	err := u.srv.Users().Delete(c, username)
	httputil.WriteResponse(c, err, nil)
}

// DeleteCollection batch delete users by usernames.
func (u *UserController) DeleteCollection(c *gin.Context) {
	u.log.X(c).Info("user delete-collection")
	usernames := c.QueryArray("name")
	u.log.Debugf("user delete-collection params: %v", usernames)

	err := u.srv.Users().DeleteCollection(c, usernames)
	httputil.WriteResponse(c, err, nil)
}
