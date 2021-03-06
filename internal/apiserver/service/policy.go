package service

import (
	"context"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"

	"github.com/che-kwas/iam-kit/meta"
)

// PolicySrv defines functions used to handle policy request.
type PolicySrv interface {
	Create(ctx context.Context, username string, policy *v1.Policy) error
	Get(ctx context.Context, username string, name string) (*v1.Policy, error)
	Update(ctx context.Context, username string, name string, params *v1.Policy) error
	List(ctx context.Context, username string, opts meta.ListOptions) (*v1.PolicyList, error)
	Delete(ctx context.Context, username string, name string) error
	DeleteCollection(ctx context.Context, username string, names []string) error
}

type policyService struct {
	store store.Store
}

var _ PolicySrv = &policyService{}

func newPolicies(srv *service) *policyService {
	return &policyService{store: srv.store}
}

func (s *policyService) Create(ctx context.Context, username string, policy *v1.Policy) error {
	policy.Username = username
	if policy.Policy.ID == "" {
		policy.Policy.ID = policy.Name
	}

	return s.store.Policies().Create(ctx, policy)
}

func (s *policyService) Get(ctx context.Context, username, name string) (*v1.Policy, error) {
	return s.store.Policies().Get(ctx, username, name)
}

func (s *policyService) Update(ctx context.Context, username string, name string, params *v1.Policy) error {
	policy, err := s.store.Policies().Get(ctx, username, name)
	if err != nil {
		return err
	}

	policy.Policy = params.Policy
	policy.Extend = params.Extend
	if policy.Policy.ID == "" {
		policy.Policy.ID = policy.Name
	}

	return s.store.Policies().Update(ctx, policy)
}

func (s *policyService) List(ctx context.Context, username string, opts meta.ListOptions) (*v1.PolicyList, error) {
	return s.store.Policies().List(ctx, username, opts)
}

func (s *policyService) Delete(ctx context.Context, username, name string) error {
	return s.store.Policies().Delete(ctx, username, name)
}

func (s *policyService) DeleteCollection(ctx context.Context, username string, names []string) error {
	return s.store.Policies().DeleteCollection(ctx, username, names)
}
