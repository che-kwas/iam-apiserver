package v1

import (
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"github.com/che-kwas/iam-kit/util"
	"gorm.io/gorm"
)

// Secret represents a secret resource.
// It is also used as gorm model.
type Secret struct {
	// Standard object metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Username          string `json:"username"    gorm:"column:username"    validate:"omitempty"`
	SecretID          string `json:"secretID"    gorm:"column:secretID"    validate:"omitempty"`
	SecretKey         string `json:"secretKey"   gorm:"column:secretKey"   validate:"omitempty"`
	Expires           int64  `json:"expires"     gorm:"column:expires"     validate:"omitempty"`
}

// SecretList is the whole list of all secrets which have been stored in stroage.
type SecretList struct {
	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	// List of secrets
	Items []*Secret `json:"items"`
}

// AfterCreate runs after create database record.
func (s *Secret) AfterCreate(tx *gorm.DB) error {
	s.InstanceID = util.GetInstanceID(s.ID, "secret-")

	return tx.Save(s).Error
}
