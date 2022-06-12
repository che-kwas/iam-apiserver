// Package mysql implements `iam-apiserver/internal/apiserver/store.Store` interface.
package mysql

import (
	"sync"

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

var (
	mysqlStore store.Store
	once       sync.Once
)

// MySQLStore returns a mysql store instance.
func MySQLStore() (store.Store, error) {
	if mysqlStore != nil {
		return mysqlStore, nil
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = mysql.NewMysqlIns()
		mysqlStore = &datastore{dbIns}
	})

	if err != nil {
		return nil, err
	}

	autoMigrate(dbIns)

	return mysqlStore, nil
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
