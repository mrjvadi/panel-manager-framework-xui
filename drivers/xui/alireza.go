package xui

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const AlirezaName = "xui.alireza"

func init() { core.Register(AlirezaName, NewAlireza) }

type alireza struct{ *generic }

func NewAlireza(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
	g := newGeneric(sp, opts...)
	return &alireza{g}, nil
}

func (d *alireza) Name() string { return AlirezaName }
func (d *alireza) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	return "alireza"
}
func (d *alireza) Capabilities() core.Capabilities                          { return core.Capabilities{} }
func (d *alireza) ListUsers(ctx context.Context) ([]core.User, error)       { return nil, nil }
func (d *alireza) ListInbounds(ctx context.Context) ([]core.Inbound, error) { return nil, nil }

var _ ext.XUITyped = (*alireza)(nil)

func (d *alireza) GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error) {
	s := &sanaei{d.generic}
	return s.GetInboundTyped(ctx, inboundID)
}
func (d *alireza) AddInboundTyped(ctx context.Context, in xdto.InboundCreate) (xdto.Inbound, error) {
	s := &sanaei{d.generic}
	return s.AddInboundTyped(ctx, in)
}
func (d *alireza) UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error) {
	s := &sanaei{d.generic}
	return s.UpdateInboundTyped(ctx, inboundID, in)
}
func (d *alireza) ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error) {
	return nil, nil
}
func (d *alireza) CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
	s := &sanaei{d.generic}
	return s.CloneInboundTyped(ctx, inboundID, opts)
}

// اگر ep قبلاً نداری، مثل Sanaei اضافه کن
func (d *alireza) ep(key, def string) string {
	if s := d.sp.Endpoints[key]; s != "" {
		return s
	}
	return def
}

func (d *alireza) DeleteInboundByID(ctx context.Context, id int) error {
	path := fmt.Sprintf(d.ep("deleteInbound", "/panel/api/inbounds/del/%d"), id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, nil)
	return d.doJSON(ctx, req, nil)
}
