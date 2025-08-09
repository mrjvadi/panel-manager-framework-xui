package ext

import (
    "context"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
)

// XUITyped: نسخه‌ی type-safe برای خانواده‌ی X-UI
type XUITyped interface {
    GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error)
    AddInboundTyped(ctx context.Context, in xdto.InboundCreate) (xdto.Inbound, error)
    UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error)
    ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error)
    CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error)
}
