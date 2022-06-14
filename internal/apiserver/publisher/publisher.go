// Package publisher defines the Publisher interface.
package publisher

import "context"

//go:generate mockgen -self_package=iam-apiserver/internal/apiserver/publisher -destination mock_publisher.go -package publisher iam-apiserver/internal/apiserver/publisher Publisher

var pub Publisher

// Publisher defines the behavior of a publisher.
type Publisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Close(ctx context.Context) error
}

// Pub returns the publisher pub.
func Pub() Publisher {
	return pub
}

// SetPub sets the publisher pub.
func SetPub(publisher Publisher) {
	pub = publisher
}
