// Package store defines the Store interface.
package store

import "context"

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/store -destination mock_store.go -package store iam-apiserver/internal/apiserver/store Store,UserStore,SecretStore,PolicyStore

var client Store

// Store defines the storage interface.
type Store interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
	Close(ctx context.Context) error
}

// Client returns the store client.
func Client() Store {
	return client
}

// SetClient sets the store client.
func SetClient(store Store) {
	client = store
}
