package mysql

import (
	"context"
	"regexp"

	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/db"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/pkg/code"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user.
func (u *users) Create(ctx context.Context, user *v1.User) error {
	if err := u.db.Create(&user).Error; err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}

		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// Get return an user by the username.
func (u *users) Get(ctx context.Context, username string) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("username = ? and isActive", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return user, nil
}

// Update updates the user.
func (u *users) Update(ctx context.Context, user *v1.User) error {
	if err := u.db.Save(user).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil

}

// List return users.
func (u *users) List(ctx context.Context, opts meta.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := db.NewOffsetLimit(opts.Offset, opts.Limit)

	err := u.db.Where("isActive").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error
	if err != nil {
		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return ret, nil
}

// Delete deletes the user by the username.
func (u *users) Delete(ctx context.Context, username string) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteByUser(ctx, username); err != nil {
		return err
	}

	if err := u.db.Where("username = ?", username).Delete(&v1.User{}).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteCollectionByUser(ctx, usernames); err != nil {
		return err
	}

	if err := u.db.Where("username in (?)", usernames).Delete(&v1.User{}).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}
