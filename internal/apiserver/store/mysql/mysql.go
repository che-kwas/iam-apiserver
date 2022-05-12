package mysql

import (
	"fmt"
	"sync"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/store"

	"github.com/che-kwas/iam-kit/db"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
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
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLFactory gets/creats a mysql factory.
func GetMySQLFactory() (store.Factory, error) {
	if mysqlFactory != nil {
		return mysqlFactory, nil
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = db.NewDBBuilder().Build()
		mysqlFactory = &datastore{dbIns}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory: %v", err)
	}

	migrateDatabase(dbIns)
	return mysqlFactory, nil
}

func migrateDatabase(db *gorm.DB) error {
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
