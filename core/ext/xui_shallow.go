package ext

import (
	"context"

	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
)

// کلون خام: فقط remark/port را عوض می‌کند؛ اگر opts.Client != nil بود، همان لحظه کلاینت هم اضافه می‌کند.
type XUIShallowClone interface {
	CloneInboundShallow(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error)
}
