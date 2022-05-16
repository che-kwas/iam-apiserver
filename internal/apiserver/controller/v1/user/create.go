package user

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/auth"
	"github.com/marmotedu/errors"

	v1 "iam-apiserver/api/apiserver/v1"
)

// Create creates a new user.
func (u *UserController) Create(c *gin.Context) {
	var user v1.User

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)

		return
	}

	if errs := user.Validate(); len(errs) != 0 {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParamsErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	user.Password, _ = auth.Encrypt(user.Password)
	user.Status = 1

	// Insert the user to the storage.
	if err := u.srv.Users().Create(c, &user, metav1.CreateOptions{}); err != nil {
		httputil.WriteResponse(c, err, nil)

		return
	}

	httputil.WriteResponse(c, nil, user)
}
