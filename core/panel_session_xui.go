package core

import (
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

// CloneInbound: شورتکات بدون نیاز به ctx بیرونی
func (x PanelXUI) CloneInbound(inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
    if xt, ok := x.p.m.As[ext.XUITyped](x.p.id); ok {
        ctx, cancel := x.p.req.derive(); defer cancel()
        return xt.CloneInboundTyped(ctx, inboundID, opts)
    }
    return xdto.Inbound{}, ErrExtNotSupported
}

// کمکی‌ها
func (x PanelXUI) CloneInboundWithPort(inboundID int, port int) (xdto.Inbound, error) {
    return x.CloneInbound(inboundID, xdto.CloneInboundOptions{ Port: &port })
}
func (x PanelXUI) CloneInboundWithRemark(inboundID int, remark string) (xdto.Inbound, error) {
    return x.CloneInbound(inboundID, xdto.CloneInboundOptions{ Remark: &remark })
}
