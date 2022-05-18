package service

import (
	"context"
	"sync"

	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/marmotedu/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *v1.User) error
	Get(ctx context.Context, username string) (*v1.User, error)
	Update(ctx context.Context, username string, params *v1.User) error
	List(ctx context.Context, opts meta.ListOptions) (*v1.UserList, error)
	Delete(ctx context.Context, username string) error
	DeleteCollection(ctx context.Context, usernames []string) error
	ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error
}

type userService struct {
	store store.Store
}

var _ UserSrv = &userService{}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User) error {
	user.Password, _ = hashPassword(user.Password)
	user.Name = user.Username
	user.IsActive = true

	return u.store.Users().Create(ctx, user)
}

func (u *userService) Get(ctx context.Context, username string) (*v1.User, error) {
	return u.store.Users().Get(ctx, username)
}

func (u *userService) Update(ctx context.Context, username string, params *v1.User) error {
	user, err := u.store.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	user.Email = params.Email
	user.Phone = params.Phone
	user.Extend = params.Extend
	return u.store.Users().Update(ctx, user)
}

func (u *userService) List(ctx context.Context, opts meta.ListOptions) (*v1.UserList, error) {
	users, err := u.store.Users().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Load TotalPolicy of users
	eg, ctx := errgroup.WithContext(ctx)
	var m sync.Map

	for _, user := range users.Items {
		user := user
		eg.Go(func() error {
			policies, err := u.store.Policies().List(ctx, user.Username, meta.ListOptions{})
			if err != nil {
				return err
			}

			m.Store(user.ID, &v1.User{
				ObjectMeta: meta.ObjectMeta{
					ID:         user.ID,
					InstanceID: user.InstanceID,
					Name:       user.Name,
					Extend:     user.Extend,
					CreatedAt:  user.CreatedAt,
					UpdatedAt:  user.UpdatedAt,
				},
				Username:    user.Username,
				Email:       user.Email,
				Phone:       user.Phone,
				TotalPolicy: policies.TotalCount,
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	infos := make([]*v1.User, 0, len(users.Items))
	for _, user := range users.Items {
		info, _ := m.Load(user.ID)
		infos = append(infos, info.(*v1.User))
	}

	return &v1.UserList{ListMeta: users.ListMeta, Items: infos}, nil
}

func (u *userService) Delete(ctx context.Context, username string) error {
	return u.store.Users().Delete(ctx, username)
}

func (u *userService) DeleteCollection(ctx context.Context, usernames []string) error {
	return u.store.Users().DeleteCollection(ctx, usernames)
}

func (u *userService) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	user, err := u.store.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if !user.VerifyPassword(oldPassword) {
		return errors.WithCode(basecode.ErrPermissionDenied, "Old password error")
	}

	user.Password, _ = hashPassword(newPassword)

	return u.store.Users().Update(ctx, user)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}
