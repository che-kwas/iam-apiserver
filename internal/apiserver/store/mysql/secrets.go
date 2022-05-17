package mysql

import (
	"context"

	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/db"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/pkg/code"
)

type secrets struct {
	db *gorm.DB
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds.db}
}

// Create creates a new secret.
func (s *secrets) Create(ctx context.Context, secret *v1.Secret) error {
	if err := s.db.Create(&secret).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// Get returns the secret by the secret identifier.
func (s *secrets) Get(ctx context.Context, username, name string) (*v1.Secret, error) {
	secret := &v1.Secret{}
	err := s.db.Where("username = ? and name= ?", username, name).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrSecretNotFound, err.Error())
		}

		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return secret, nil
}

// Update updates the secret.
func (s *secrets) Update(ctx context.Context, secret *v1.Secret) error {
	if err := s.db.Save(secret).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// List returns all secrets.
func (s *secrets) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	ret := &v1.SecretList{}
	ol := db.NewOffsetLimit(opts.Offset, opts.Limit)

	if username != "" {
		s.db = s.db.Where("username = ?", username)
	}

	err := s.db.
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

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(ctx context.Context, username, name string) error {
	err := s.db.Where("username = ? and name = ?", username, name).Delete(&v1.Secret{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the secrets.
func (s *secrets) DeleteCollection(ctx context.Context, username string, names []string) error {
	err := s.db.Where("username = ? and name in (?)", username, names).Delete(&v1.Secret{}).Error
	if err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}
