package core

import "context"

type Driver interface {
	Name() string
	Version() string
	Capabilities() Capabilities
	Login(ctx context.Context) error
	ListUsers(ctx context.Context) ([]User, error)
	ListInbounds(ctx context.Context) ([]Inbound, error)
}
