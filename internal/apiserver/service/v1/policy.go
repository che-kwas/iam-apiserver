package v1

import (
	"context"

	metav1 "github.com/che-kwas/iam-kit/meta/v1"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

// PolicySrv defines functions used to handle policy request.
type PolicySrv interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error
	Get(ctx context.Context, username string, name string, opts metav1.GetOptions) (*v1.Policy, error)
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error)
	Delete(ctx context.Context, username string, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOptions) error
}

type policyService struct {
	store store.Store
}

var _ PolicySrv = &policyService{}

func newPolicies(srv *service) *policyService {
	return &policyService{store: srv.store}
}

func (s *policyService) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	return s.store.Policies().Create(ctx, policy, opts)
}

func (s *policyService) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	return s.store.Policies().Get(ctx, username, name, opts)
}

func (s *policyService) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	return s.store.Policies().Update(ctx, policy, opts)
}

func (s *policyService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	return s.store.Policies().List(ctx, username, opts)
}

func (s *policyService) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	return s.store.Policies().Delete(ctx, username, name, opts)
}

func (s *policyService) DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOptions) error {
	return s.store.Policies().DeleteCollection(ctx, username, names, opts)
}
