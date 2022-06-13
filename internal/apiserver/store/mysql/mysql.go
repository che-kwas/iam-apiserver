// Package mysql implements `iam-apiserver/internal/apiserver/store.Store` interface.
package mysql

import (
	"github.com/che-kwas/iam-kit/mysql"
	"gorm.io/gorm"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

// NewMySQLStore returns a mysql store instance.
func NewMySQLStore() (store.Store, error) {
	dbIns, err := mysql.NewMysqlIns()
	if err != nil {
		return nil, err
	}

	autoMigrate(dbIns)

	return &datastore{dbIns}, nil
}

// nolint:unused
func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&v1.Policy{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&v1.Secret{}); err != nil {
		return err
	}

	return nil
}
