// Package publisher defines the Publisher interface.
package publisher

import "context"

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/publisher -destination mock_publisher.go -package publisher iam-apiserver/internal/apiserver/publisher Publisher

var client Publisher

// Publisher defines the behavior of a publisher.
type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Close() error
}

// Client returns the publisher client.
func Client() Publisher {
	return client
}

// SetClient sets the publisher client.
func SetClient(publisher Publisher) {
	client = publisher
}
