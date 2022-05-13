package mysql

import (
	"context"

	v1 "iam-apiserver/api/apiserver/v1"

	"github.com/che-kwas/iam-kit/db"
	"github.com/che-kwas/iam-kit/errcode"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

// Create creates a new user.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.db.Create(&user).Error
}

// Update updates the user.
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the username.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteByUser(ctx, username, opts); err != nil {
		return err
	}

	err := u.db.Where("username = ?", username).Delete(&v1.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(errcode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	// delete related policy first
	pol := newPolicies(&datastore{u.db})
	if err := pol.DeleteCollectionByUser(ctx, usernames, opts); err != nil {
		return err
	}

	return u.db.Where("name in (?)", usernames).Delete(&v1.User{}).Error
}

// Get return an user by the username.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("username = ? and isActive", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrNotFound, err.Error())
		}

		return nil, errors.WithCode(errcode.ErrDatabase, err.Error())
	}

	return user, nil
}

// List return users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := db.NewOffsetLimit(opts.Offset, opts.Limit)

	d := u.db.Where("isActive").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
