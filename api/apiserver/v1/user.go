package v1

import (
	"time"

	"github.com/che-kwas/iam-kit/meta"
	"github.com/che-kwas/iam-kit/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user resource.
// It is also used as gorm model.
type User struct {
	// Standard object metadata.
	meta.ObjectMeta `json:"metadata,omitempty"`

	Username    string     `json:"username"              gorm:"column:username"   validate:"required,min=1,max=30"`
	Password    string     `json:"password,omitempty"    gorm:"column:password"   validate:"required"`
	Email       string     `json:"email"                 gorm:"column:email"      validate:"required,email,min=1,max=100"`
	Phone       string     `json:"phone,omitempty"       gorm:"column:phone"`
	TotalPolicy int64      `json:"totalPolicy,omitempty" gorm:"-"`
	LoginedAt   *time.Time `json:"loginedAt,omitempty"   gorm:"column:loginedAt"`
	IsAdmin     bool       `json:"isAdmin,omitempty"     gorm:"column:isAdmin;default:false"`
	IsActive    bool       `json:"isActive,omitempty"    gorm:"column:isActive;default:true"`
}

// AfterCreate runs after create database record.
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = util.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}

// VerifyPassword verifies the plain text password.
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserList represents a collection of users.
type UserList struct {
	// Standard list metadata.
	meta.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}
