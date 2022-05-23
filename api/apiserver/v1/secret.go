package v1

import (
	"github.com/che-kwas/iam-kit/meta"
	"github.com/che-kwas/iam-kit/util"
	"gorm.io/gorm"
)

// Secret represents a secret resource.
// It is also used as gorm model.
type Secret struct {
	// Standard object metadata.
	meta.ObjectMeta `json:"metadata,omitempty"`
	Username        string `json:"username,omitempty"  gorm:"column:username"`
	SecretID        string `json:"secretID,omitempty"  gorm:"column:secretID"`
	SecretKey       string `json:"secretKey,omitempty" gorm:"column:secretKey"`
	Expires         int64  `json:"expires"             gorm:"column:expires"     validate:"required"`
	Description     string `json:"description"         gorm:"column:description" validate:"required,email,min=1,max=100"`
}

// SecretList is the whole list of all secrets which have been stored in stroage.
type SecretList struct {
	// Standard list metadata.
	meta.ListMeta `json:",inline"`

	// List of secrets
	Items []*Secret `json:"items"`
}

// AfterCreate runs after create database record.
func (s *Secret) AfterCreate(tx *gorm.DB) error {
	s.InstanceID = util.GetInstanceID(s.ID, "secret-")

	return tx.Save(s).Error
}
