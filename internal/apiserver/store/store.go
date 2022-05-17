// Package store defines the storage interface for iam-apiserver.
package store

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/store -destination mock_store.go -package store iam-apiserver/internal/apiserver/store Factory,UserStore,SecretStore,PolicyStore

var client Store

// Store defines the apiserver storage interface.
type Store interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
	Close() error
}

// Client return the store client.
func Client() Store {
	return client
}

// SetClient set the store client.
func SetClient(store Store) {
	client = store
}
