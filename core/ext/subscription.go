package ext

import "context"

// Subscription: عملیات سابسکریپشن (Marzban)

type Subscription interface {
    CreateSubscription(ctx context.Context, payload map[string]any) (map[string]any, error)
    GetSubscription(ctx context.Context, id string) (map[string]any, error)
    ListSubscriptions(ctx context.Context) ([]map[string]any, error)
    DeleteSubscription(ctx context.Context, id string) error
    RevokeUserSubscription(ctx context.Context, username string) error
}
