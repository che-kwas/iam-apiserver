package service

import (
	"context"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/che-kwas/iam-kit/util"
	"github.com/marmotedu/errors"
	"github.com/spf13/viper"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
	"iam-apiserver/internal/pkg/code"
	"iam-apiserver/internal/pkg/token"
)

const (
	secretIDLen    = 36
	secretKeyLen   = 32
	maxSecretCount = 10
)

// SecretSrv defines functions used to handle secret request.
type SecretSrv interface {
	Create(ctx context.Context, username string, secret *v1.Secret) error
	Get(ctx context.Context, username, name string) (*v1.Secret, error)
	GetToken(ctx context.Context, username, name string) (*v1.Token, error)
	Update(ctx context.Context, username, name string, params *v1.Secret) error
	List(ctx context.Context, username string, opts meta.ListOptions) (*v1.SecretList, error)
	Delete(ctx context.Context, username, name string) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string) error
}

type secretService struct {
	store store.Store
}

var _ SecretSrv = &secretService{}

func newSecrets(srv *service) *secretService {
	return &secretService{store: srv.store}
}

func (s *secretService) Create(ctx context.Context, username string, secret *v1.Secret) error {
	// check TotalCount
	listOpts := meta.ListOptions{Offset: pointer.ToInt(0), Limit: pointer.ToInt(-1)}
	secrets, err := s.store.Secrets().List(ctx, username, listOpts)
	if err != nil {
		return err
	}
	if secrets.TotalCount >= maxSecretCount {
		return errors.WithCode(code.ErrReachMaxCount, "secret count: %d", secrets.TotalCount)
	}

	secret.Username = username
	secret.SecretID = util.RandString(secretIDLen)
	secret.SecretKey = util.RandString(secretKeyLen)
	if secret.Name == "" {
		secret.Name = secret.SecretID
	}

	return s.store.Secrets().Create(ctx, secret)
}

func (s *secretService) Get(ctx context.Context, username, name string) (*v1.Secret, error) {
	return s.store.Secrets().Get(ctx, username, name)
}

func (s *secretService) GetToken(ctx context.Context, username, name string) (*v1.Token, error) {
	secret, err := s.store.Secrets().Get(ctx, username, name)
	if err != nil {
		return nil, err
	}

	timeout := viper.GetDuration("jwt.timeout")
	expire := time.Now().Add(timeout)
	token, _ := token.GenToken(secret.SecretID, secret.SecretKey, expire)

	return &v1.Token{Token: token, Expire: expire.Format(time.RFC3339)}, nil
}

func (s *secretService) Update(ctx context.Context, username, name string, params *v1.Secret) error {
	secret, err := s.store.Secrets().Get(ctx, username, name)
	if err != nil {
		return err
	}

	secret.Expires = params.Expires
	secret.Description = params.Description
	secret.Extend = params.Extend

	return s.store.Secrets().Update(ctx, secret)
}

func (s *secretService) List(ctx context.Context, username string, opts meta.ListOptions) (*v1.SecretList, error) {
	return s.store.Secrets().List(ctx, username, opts)
}

func (s *secretService) Delete(ctx context.Context, username, name string) error {
	return s.store.Secrets().Delete(ctx, username, name)
}

func (s *secretService) DeleteCollection(ctx context.Context, username string, secretIDs []string) error {
	return s.store.Secrets().DeleteCollection(ctx, username, secretIDs)
}
