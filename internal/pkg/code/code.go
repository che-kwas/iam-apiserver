// Package code defines error codes for iam apiserver.
package code

//go:generate codegen
//go:generate codegen -doc -output ../../../errcode_apiserver.md

// user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist
)

// secret errors.
const (
	// ErrSecretNotFound - 404: Secret not found.
	ErrSecretNotFound int = iota + 110101

	// ErrEncrypt - 400: Secret reach the max count.
	ErrReachMaxCount
)

// policy errors.
const (
	// ErrPolicyNotFound - 404: Policy not found.
	ErrPolicyNotFound int = iota + 110201
)
