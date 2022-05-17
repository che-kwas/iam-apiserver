package store

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
)

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	Create(ctx context.Context, policy *v1.Policy) error
	Get(ctx context.Context, username string, name string) (*v1.Policy, error)
	Update(ctx context.Context, policy *v1.Policy) error
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error)
	Delete(ctx context.Context, username string, name string) error
	DeleteCollection(ctx context.Context, username string, names []string) error
}
