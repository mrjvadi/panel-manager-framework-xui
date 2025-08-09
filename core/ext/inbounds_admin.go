package ext

import "context"

// InboundsAdmin: عملیات خاص X-UI روی این‌باند و کلاینت‌ها

type InboundsAdmin interface {
    GetInbound(ctx context.Context, inboundID int) (map[string]any, error)
    AddInbound(ctx context.Context, payload map[string]any) (map[string]any, error)
    UpdateInbound(ctx context.Context, inboundID int, payload map[string]any) (map[string]any, error)
    DeleteInbound(ctx context.Context, inboundID int) error
    AddClient(ctx context.Context, payload map[string]any) (map[string]any, error)
    DeleteClient(ctx context.Context, clientID int) error
    ClientTrafficByEmail(ctx context.Context, email string) (map[string]any, error)
    ClientTrafficByID(ctx context.Context, uuid string) (map[string]any, error)
    ResetAllTraffic(ctx context.Context) error
    ClientIPs(ctx context.Context, email string) ([]string, error)
}
