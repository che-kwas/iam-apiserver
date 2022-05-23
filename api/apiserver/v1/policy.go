package v1

import (
	"encoding/json"

	"github.com/che-kwas/iam-kit/meta"
	"github.com/che-kwas/iam-kit/util"
	"github.com/ory/ladon"
	"gorm.io/gorm"
)

// Policy represents a policy resource including a ladon policy.
// It is also used as gorm model.
type Policy struct {
	// Standard object metadata.
	meta.ObjectMeta `json:"metadata,omitempty"`

	// The user of the policy.
	Username string `json:"username,omitempty" gorm:"column:username"`

	Policy ladon.DefaultPolicy `json:"policy,omitempty" gorm:"-"`

	// The ladon policy content, just a string format of ladon.DefaultPolicy. DO NOT modify directly.
	PolicyShadow string `json:"-" gorm:"column:policyShadow" validate:"omitempty"`
}

// PolicyList is the whole list of all policies which have been stored in stroage.
type PolicyList struct {
	// Standard list metadata.
	meta.ListMeta `json:",inline"`

	// List of policies.
	Items []*Policy `json:"items"`
}

// BeforeCreate runs before create database record.
func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeCreate(tx); err != nil {
		return err
	}

	data, err := json.Marshal(p.Policy)
	if err != nil {
		return err
	}

	p.PolicyShadow = string(data)
	return nil
}

// AfterCreate runs after create database record.
func (p *Policy) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = util.GetInstanceID(p.ID, "policy-")

	return tx.Save(p).Error
}

// BeforeUpdate runs before update database record.
func (p *Policy) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeUpdate(tx); err != nil {
		return err
	}

	data, err := json.Marshal(p.Policy)
	if err != nil {
		return err
	}

	p.PolicyShadow = string(data)
	return nil
}

// AfterFind runs after find to unmarshal a policy string into ladon.DefaultPolicy struct.
func (p *Policy) AfterFind(tx *gorm.DB) (err error) {
	if err := p.ObjectMeta.AfterFind(tx); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(p.PolicyShadow), &p.Policy); err != nil {
		return err
	}

	return nil
}
