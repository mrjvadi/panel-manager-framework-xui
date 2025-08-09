package xui

import (
    "context"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    extpkg "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const AlirezaName = "xui.alireza"

func init() { core.Register(AlirezaName, NewAlireza) }

type alireza struct { *generic }

func NewAlireza(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    if sp.Endpoints == nil { sp.Endpoints = map[string]string{} }
    if sp.Endpoints["login"] == "" { sp.Endpoints["login"] = "/api/auth/login" }
    if sp.Endpoints["listUsers"] == "" { sp.Endpoints["listUsers"] = "/api/users" }
    if sp.Endpoints["listInbounds"] == "" { sp.Endpoints["listInbounds"] = "/xui/inbound/list" }
    if sp.Version == "" { sp.Version = "alireza" }
    g := newGeneric(sp, opts...)
    return &alireza{ g }, nil
}

func (d *alireza) Name() string    { return AlirezaName }
func (d *alireza) Version() string { if d.sp.Version != "" { return d.sp.Version }; return "alireza" }

// Typed API conformance via Sanaei-compatible endpoints
var _ extpkg.XUITyped = (*alireza)(nil)

func (d *alireza) GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error) {
    s := &sanaei{ d.generic }
    return s.GetInboundTyped(ctx, inboundID)
}
func (d *alireza) AddInboundTyped(ctx context.Context, in xdto.InboundCreate) (xdto.Inbound, error) {
    s := &sanaei{ d.generic }
    return s.AddInboundTyped(ctx, in)
}
func (d *alireza) UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error) {
    s := &sanaei{ d.generic }
    return s.UpdateInboundTyped(ctx, inboundID, in)
}
func (d *alireza) ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error) {
    s := &sanaei{ d.generic }
    return s.ClientTrafficByEmailTyped(ctx, email)
}
func (d *alireza) CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
    s := &sanaei{ d.generic }
    return s.CloneInboundTyped(ctx, inboundID, opts)
}
