// Package cache is the cache controller which can return all secrets and policies.
package cache

import (
	"context"

	"github.com/AlekSi/pointer"
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/marmotedu/errors"

	pb "iam-apiserver/api/apiserver/proto/v1"
	"iam-apiserver/internal/apiserver/service"
)

// CacheController handles requests for listing all secrets and policies.
type CacheController struct {
	pb.UnimplementedCacheServer

	srv service.Service
	log *logger.Logger
}

var _ pb.CacheServer = &CacheController{}

// NewCacheController creates a secret controller.
func NewCacheController() *CacheController {
	return &CacheController{
		srv: service.NewService(),
		log: logger.L(),
	}
}

// ListSecrets returns all secrets.
func (c *CacheController) ListSecrets(ctx context.Context, r *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	c.log.X(ctx).Info("secret list")

	opts := meta.ListOptions{
		Offset: pointer.ToInt(int(r.Offset)),
		Limit:  pointer.ToInt(int(r.Limit)),
	}
	c.log.Debugf("secret list params: %+v", opts)

	secrets, err := c.srv.Secrets().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	items := make([]*pb.SecretInfo, 0, len(secrets.Items))
	for _, secret := range secrets.Items {
		items = append(items, &pb.SecretInfo{
			SecretId:    secret.SecretID,
			Username:    secret.Username,
			SecretKey:   secret.SecretKey,
			Expires:     secret.Expires,
			Description: secret.Description,
			CreatedAt:   secret.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   secret.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListSecretsResponse{
		TotalCount: secrets.TotalCount,
		Items:      items,
	}, nil
}

// ListPolicies returns all policies.
func (c *CacheController) ListPolicies(ctx context.Context, r *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	c.log.X(ctx).Info("policy list")

	opts := meta.ListOptions{
		Offset: pointer.ToInt(int(r.Offset)),
		Limit:  pointer.ToInt(int(r.Limit)),
	}
	c.log.Debugf("policy list params: %+v", opts)

	policies, err := c.srv.Policies().List(ctx, "", opts)
	if err != nil {
		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	items := make([]*pb.PolicyInfo, 0, len(policies.Items))
	for _, pol := range policies.Items {
		items = append(items, &pb.PolicyInfo{
			Name:         pol.Name,
			Username:     pol.Username,
			PolicyShadow: pol.PolicyShadow,
			CreatedAt:    pol.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListPoliciesResponse{
		TotalCount: policies.TotalCount,
		Items:      items,
	}, nil
}
