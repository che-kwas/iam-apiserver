package store

import (
	"context"

	v1 "iam-apiserver/api/apiserver/v1"

	"github.com/che-kwas/iam-kit/meta"
)

// SecretStore defines the secret storage interface.
type SecretStore interface {
	Create(ctx context.Context, secret *v1.Secret) error
	Get(ctx context.Context, username, secretID string) (*v1.Secret, error)
	Update(ctx context.Context, secret *v1.Secret) error
	List(ctx context.Context, username string, opts meta.ListOptions) (*v1.SecretList, error)
	Delete(ctx context.Context, username, secretID string) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string) error
}
