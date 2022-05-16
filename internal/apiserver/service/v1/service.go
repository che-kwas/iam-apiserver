package v1

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/service/v1 -destination mock_service.go -package v1 iam-apiserver/internal/apiserver/service/v1 Service,UserSrv,SecretSrv,PolicySrv

import "iam-apiserver/internal/apiserver/store"

// Service defines functions used to return resource interface.
type Service interface {
	Users() UserSrv
	Secrets() SecretSrv
	Policies() PolicySrv
}

type service struct {
	store store.Store
}

// NewService returns Service interface.
func NewService() Service {
	return &service{
		store: store.Client(),
	}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func (s *service) Secrets() SecretSrv {
	return newSecrets(s)
}

func (s *service) Policies() PolicySrv {
	return newPolicies(s)
}
