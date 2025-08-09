package core

import (
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

type PanelXUI struct{ p *Panel }

func (p *Panel) XUI() PanelXUI { return PanelXUI{p: p} }

// شورتکات: ctx داخلی می‌سازد و کلون خام را صدا می‌زند
func (x PanelXUI) CloneInboundShallow(inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
	if sc, ok := As[ext.XUIShallowClone](x.p.m, x.p.id); ok {
		ctx, cancel := x.p.req.derive()
		defer cancel()
		return sc.CloneInboundShallow(ctx, inboundID, opts)
	}
	return xdto.Inbound{}, ErrExtNotSupported
}
