package core

import "context"

// Driver: قرارداد سراسری برای پنل‌ها

type Driver interface {
    Name() string
    Version() string
    Capabilities() Capabilities
    Login(ctx context.Context) error
    ListUsers(ctx context.Context) ([]User, error)
    ListInbounds(ctx context.Context) ([]Inbound, error)
    CreateUser(ctx context.Context, u User) (User, error)
    UpdateUser(ctx context.Context, u User) (User, error)
    DeleteUser(ctx context.Context, id string) error
    SuspendUser(ctx context.Context, id string) error
    ResumeUser(ctx context.Context, id string) error
    ResetUserTraffic(ctx context.Context, id string) error
    CreateInbound(ctx context.Context, in Inbound) (Inbound, error)
    UpdateInbound(ctx context.Context, in Inbound) (Inbound, error)
    DeleteInbound(ctx context.Context, id string) error
}

// DriverFactory برای رجیستری

type DriverFactory func(PanelSpec, ...Option) (Driver, error)

// PanelSpec تعریف کُد-محور پنل

type PanelSpec struct {
    ID        string
    BaseURL   string
    Auth      Auth // Basic/Token/NoAuth
    TLS       TLS
    Version   string            // اختیاری: override نسخه‌ی درایور/پنل
    Endpoints map[string]string // override مسیرها
}
