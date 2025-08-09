package ext

import "context"

type Usage interface {
    UserUsage(ctx context.Context, username string) (map[string]any, error)
    UsersUsage(ctx context.Context) ([]map[string]any, error)
    UsersExpired(ctx context.Context) ([]map[string]any, error)
    ResetAllUsers(ctx context.Context) error
}
