package v1

import (
	"time"

	"github.com/che-kwas/iam-kit/db"
	metav1 "github.com/che-kwas/iam-kit/meta/v1"
	"gorm.io/gorm"
)

// User represents a user resource.
// It is also used as gorm model.
type User struct {
	// Standard object metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Username    string    `json:"username"            gorm:"column:username"   validate:"required,min=1,max=30"`
	Password    string    `json:"password,omitempty"  gorm:"column:password"   validate:"required"`
	Email       string    `json:"email"               gorm:"column:email"      validate:"required,email,min=1,max=100"`
	Phone       string    `json:"phone"               gorm:"column:phone"      validate:"omitempty"`
	TotalPolicy int64     `json:"totalPolicy"         gorm:"-"                 validate:"omitempty"`
	IsAdmin     bool      `json:"isAdmin,omitempty"   gorm:"column:isAdmin;default:false" validate:"omitempty"`
	IsActive    bool      `json:"isActive,omitempty"  gorm:"column:isActive;default:true" validate:"omitempty"`
	LoginedAt   time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

// AfterCreate runs after create database record.
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = db.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}
