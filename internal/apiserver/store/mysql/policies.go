package mysql

import (
	"context"

	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/db"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/pkg/code"
)

type policies struct {
	db *gorm.DB
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds.db}
}

// Create creates a new ladon policy.
func (p *policies) Create(ctx context.Context, policy *v1.Policy) error {
	if err := p.db.Create(&policy).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// Get returns the policy by the policy identifier.
func (p *policies) Get(ctx context.Context, username, name string) (*v1.Policy, error) {
	policy := &v1.Policy{}
	err := p.db.Where("username = ? and name = ?", username, name).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return policy, nil
}

// Update updates the policy.
func (p *policies) Update(ctx context.Context, policy *v1.Policy) error {
	if err := p.db.Save(policy).Error; err != nil {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// List returns all policies.
func (p *policies) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	ret := &v1.PolicyList{}
	ol := db.NewOffsetLimit(opts.Offset, opts.Limit)

	if username != "" {
		p.db = p.db.Where("username = ?", username)
	}

	err := p.db.
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount).
		Error
	if err != nil {
		return nil, errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return ret, nil
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(ctx context.Context, username, name string) error {
	err := p.db.Where("username = ? and name = ?", username, name).Delete(&v1.Policy{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteByUser deletes policies by username.
func (p *policies) DeleteByUser(ctx context.Context, username string) error {
	err := p.db.Where("username = ?", username).Delete(&v1.Policy{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes policies by policies ids.
func (p *policies) DeleteCollection(ctx context.Context, username string, names []string) error {
	err := p.db.Where("username = ? and name in (?)", username, names).Delete(&v1.Policy{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(ctx context.Context, usernames []string) error {
	err := p.db.Where("username in (?)", usernames).Delete(&v1.Policy{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(basecode.ErrDatabase, err.Error())
	}

	return nil
}
