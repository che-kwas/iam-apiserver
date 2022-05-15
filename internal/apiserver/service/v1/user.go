package v1

import (
	"context"
	"regexp"

	"github.com/che-kwas/iam-kit/errcode"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/marmotedu/errors"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	// Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	// Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	// DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	// Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	// List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	// ChangePassword(ctx context.Context, user *v1.User) error
}

type userService struct {
	store store.Store
}

var _ UserSrv = &userService{}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	if err := u.store.Users().Create(ctx, user, opts); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			// TODO define api-server errors
			return errors.WithCode(errcode.ErrDatabase, err.Error())
		}

		return errors.WithCode(errcode.ErrDatabase, err.Error())
	}

	return nil
}
