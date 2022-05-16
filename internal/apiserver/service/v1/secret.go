package v1

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// SecretSrv defines functions used to handle secret request.
type SecretSrv interface {
	Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error
	Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error)
	Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
	Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOptions) error
}

type secretService struct {
	store store.Store
}

var _ SecretSrv = &secretService{}

func newSecrets(srv *service) *secretService {
	return &secretService{store: srv.store}
}

func (s *secretService) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	return s.store.Secrets().Create(ctx, secret, opts)
}

func (s *secretService) Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error) {
	return s.store.Secrets().Get(ctx, username, secretID, opts)
}

func (s *secretService) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	return s.store.Secrets().Update(ctx, secret, opts)
}

func (s *secretService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	return s.store.Secrets().List(ctx, username, opts)
}

func (s *secretService) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error {
	return s.store.Secrets().Delete(ctx, username, secretID, opts)
}

func (s *secretService) DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOptions) error {
	return s.store.Secrets().DeleteCollection(ctx, username, secretIDs, opts)
}
