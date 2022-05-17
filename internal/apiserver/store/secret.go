package store

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
)

// SecretStore defines the secret storage interface.
type SecretStore interface {
	Create(ctx context.Context, secret *v1.Secret) error
	Get(ctx context.Context, username, secretID string) (*v1.Secret, error)
	Update(ctx context.Context, secret *v1.Secret) error
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
	Delete(ctx context.Context, username, secretID string) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string) error
}
