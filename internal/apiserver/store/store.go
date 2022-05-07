package store

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/store -destination mock_store.go -package store iam-apiserver/internal/apiserver/store Factory,UserStore,SecretStore,PolicyStore

var client Factory

// Factory defines the apiserver storage interface.
type Factory interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
	Close() error
}

// Client return the store client.
func Client() Factory {
	return client
}

// SetClient set the store client.
func SetClient(factory Factory) {
	client = factory
}
