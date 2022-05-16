package user

import (
	srvv1 "iam-apiserver/internal/apiserver/service/v1"
)

// UserController handles requests for user resource.
type UserController struct {
	srv srvv1.Service
}

// NewUserController creates a user handler.
func NewUserController() *UserController {
	return &UserController{
		srv: srvv1.NewService(),
	}
}
