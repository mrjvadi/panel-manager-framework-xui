package core

import (
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
)

type PanelXUI struct{ p *Panel }

func (p *Panel) XUI() PanelXUI { return PanelXUI{ p: p } }

func (x PanelXUI) CloneInbound(inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
    if xt, ok := x.p.m.As[ext.XUITyped](x.p.id); ok {
        ctx, cancel := x.p.req.derive(); defer cancel()
        return xt.CloneInboundTyped(ctx, inboundID, opts)
    }
    return xdto.Inbound{}, ErrExtNotSupported
}
func (x PanelXUI) CloneInboundWithPort(inboundID, port int) (xdto.Inbound, error) {
    return x.CloneInbound(inboundID, xdto.CloneInboundOptions{ Port: &port })
}
func (x PanelXUI) CloneInboundWithRemark(inboundID int, remark string) (xdto.Inbound, error) {
    return x.CloneInbound(inboundID, xdto.CloneInboundOptions{ Remark: &remark })
}
