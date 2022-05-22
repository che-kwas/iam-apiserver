// Package mysql implements `iam-apiserver/internal/apiserver/store.Store` interface.
package mysql

import (
	"fmt"
	"sync"

	"github.com/che-kwas/iam-kit/db"
	"github.com/marmotedu/errors"
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
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlStore store.Store
	once       sync.Once
)

// GetMySQLStore gets/creats a mysql store.
func GetMySQLStore() (store.Store, error) {
	if mysqlStore != nil {
		return mysqlStore, nil
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = db.NewDBBuilder().Build()
		mysqlStore = &datastore{dbIns}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory: %v", err)
	}

	// TODO manual migrate
	autoMigrate(dbIns)

	return mysqlStore, nil
}

// nolint:unused
func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return errors.Wrap(err, "migrate user model failed")
	}
	if err := db.AutoMigrate(&v1.Policy{}); err != nil {
		return errors.Wrap(err, "migrate policy model failed")
	}
	if err := db.AutoMigrate(&v1.Secret{}); err != nil {
		return errors.Wrap(err, "migrate secret model failed")
	}

	return nil
}
