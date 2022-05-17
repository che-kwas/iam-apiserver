package service

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// SecretSrv defines functions used to handle secret request.
type SecretSrv interface {
	Create(ctx context.Context, secret *v1.Secret) error
	Get(ctx context.Context, username, secretID string) (*v1.Secret, error)
	Update(ctx context.Context, secret *v1.Secret) error
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
	Delete(ctx context.Context, username, secretID string) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string) error
}

type secretService struct {
	store store.Store
}

var _ SecretSrv = &secretService{}

func newSecrets(srv *service) *secretService {
	return &secretService{store: srv.store}
}

func (s *secretService) Create(ctx context.Context, secret *v1.Secret) error {
	return s.store.Secrets().Create(ctx, secret)
}

func (s *secretService) Get(ctx context.Context, username, secretID string) (*v1.Secret, error) {
	return s.store.Secrets().Get(ctx, username, secretID)
}

func (s *secretService) Update(ctx context.Context, secret *v1.Secret) error {
	return s.store.Secrets().Update(ctx, secret)
}

func (s *secretService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	return s.store.Secrets().List(ctx, username, opts)
}

func (s *secretService) Delete(ctx context.Context, username, secretID string) error {
	return s.store.Secrets().Delete(ctx, username, secretID)
}

func (s *secretService) DeleteCollection(ctx context.Context, username string, secretIDs []string) error {
	return s.store.Secrets().DeleteCollection(ctx, username, secretIDs)
}
