package store

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
)

// UserStore defines the user storage interface.
type UserStore interface {
	Create(ctx context.Context, user *v1.User) error
	Get(ctx context.Context, username string) (*v1.User, error)
	Update(ctx context.Context, user *v1.User) error
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	Delete(ctx context.Context, username string) error
	DeleteCollection(ctx context.Context, usernames []string) error
}
