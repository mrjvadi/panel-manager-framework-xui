package ext

import (
    "context"
    mdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/marzban"
)

// MarzbanTyped: نسخه‌ی type-safe از برخی متدهای مرزبان
type MarzbanTyped interface {
    SystemInfoTyped(ctx context.Context) (mdto.SystemInfo, error)
    UsersUsageTyped(ctx context.Context) ([]mdto.UserUsage, error)
    UserUsageTyped(ctx context.Context, username string) (mdto.UserUsage, error)
    ListSubscriptionsTyped(ctx context.Context) ([]mdto.Subscription, error)
}
