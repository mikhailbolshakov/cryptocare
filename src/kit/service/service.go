package service

import "context"

// Service declares an interface each service must implement
type Service interface {
	// GetCode returns the service unique code
	GetCode() string
	// Init initializes the service
	Init(ctx context.Context) error
	// ListenAsync executes all background processes
	Start(ctx context.Context) error
	// Close closes the service
	Close(ctx context.Context)
}

// StorageAdapter common interface for storage adapters
type StorageAdapter interface {
	// Init initializes adapter
	Init(ctx context.Context, cfg interface{}) error
	// Close closes storage
	Close(ctx context.Context) error
}
