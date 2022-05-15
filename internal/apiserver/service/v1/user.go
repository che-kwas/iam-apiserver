package v1

import (
	"context"
	"sync"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"golang.org/x/sync/errgroup"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// UserSrv defines functions used to handle user request.
type UserSrv interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	ChangePassword(ctx context.Context, user *v1.User) error
}

type userService struct {
	store store.Store
}

var _ UserSrv = &userService{}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.store.Users().Create(ctx, user, opts)
}

func (u *userService) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	return u.store.Users().Get(ctx, username, opts)
}

func (u *userService) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.store.Users().Update(ctx, user, opts)
}

func (u *userService) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
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
			policies, err := u.store.Policies().List(ctx, user.Username, metav1.ListOptions{})
			if err != nil {
				return err
			}

			m.Store(user.ID, &v1.User{
				ObjectMeta: metav1.ObjectMeta{
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

func (u *userService) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	return u.store.Users().Delete(ctx, username, opts)
}

func (u *userService) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	return u.store.Users().DeleteCollection(ctx, usernames, opts)
}

func (u *userService) ChangePassword(ctx context.Context, user *v1.User) error {
	// Save changed fields.
	return u.store.Users().Update(ctx, user, metav1.UpdateOptions{})
}
