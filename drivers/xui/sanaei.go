package xui

import (
    "context"
    "fmt"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const SanaeiName = "xui.sanaei"
const PluginVer = "v2.6.1"

func init() { core.Register(SanaeiName, NewSanaei) }

type sanaei struct{ *generic }

func NewSanaei(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    if sp.Endpoints == nil { sp.Endpoints = map[string]string{} }
    if sp.Endpoints["login"] == "" { sp.Endpoints["login"] = "/login" }
    if sp.Endpoints["getInbound"] == "" { sp.Endpoints["getInbound"] = "/panel/api/inbounds/get/%d" }
    if sp.Endpoints["addInbound"] == "" { sp.Endpoints["addInbound"] = "/panel/api/inbounds/add" }
    g := newGeneric(sp, opts...)
    return &sanaei{ g }, nil
}

func (d *sanaei) Name() string    { return SanaeiName }
func (d *sanaei) Version() string { if d.sp.Version != "" { return d.sp.Version }; return PluginVer }
func (d *sanaei) Capabilities() core.Capabilities { return core.Capabilities{} }
func (d *sanaei) ListUsers(ctx context.Context) ([]core.User, error) { return nil, nil }
func (d *sanaei) ListInbounds(ctx context.Context) ([]core.Inbound, error) { return nil, nil }

var _ ext.XUITyped = (*sanaei)(nil)

func (d *sanaei) GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error) {
    path := fmt.Sprintf(d.sp.Endpoints["getInbound"], inboundID)
    var m map[string]any
    if err := d.getJSON(ctx, path, &m); err != nil { return xdto.Inbound{}, err }
    out := xdto.Inbound{ Raw: m }
    if v, ok := m["id"].(float64); ok { out.ID = int(v) }
    if s, ok := m["remark"].(string); ok { out.Remark = s }
    if s, ok := m["protocol"].(string); ok { out.Protocol = s }
    if v, ok := m["port"].(float64); ok { out.Port = int(v) }
    if mm, ok := m["settings"].(map[string]any); ok { out.Settings = mm }
    if mm, ok := m["sniffing"].(map[string]any); ok { out.Sniffing = mm }
    if mm, ok := m["streamSettings"].(map[string]any); ok { out.StreamSettings = mm }
    return out, nil
}

func (d *sanaei) AddInboundTyped(ctx context.Context, in xdto.InboundCreate) (xdto.Inbound, error) {
    var m map[string]any
    if err := d.postJSON(ctx, d.sp.Endpoints["addInbound"], map[string]any{
        "remark": in.Remark, "protocol": in.Protocol, "port": in.Port,
        "settings": in.Settings, "sniffing": in.Sniffing, "streamSettings": in.StreamSettings,
    }, &m); err != nil { return xdto.Inbound{}, err }
    if v, ok := m["id"].(float64); ok && int(v) > 0 { return d.GetInboundTyped(ctx, int(v)) }
    out := xdto.Inbound{ Raw: m }
    if v, ok := m["id"].(float64); ok { out.ID = int(v) }
    if v, ok := m["port"].(float64); ok { out.Port = int(v) }
    return out, nil
}

func (d *sanaei) UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error) {
    return d.AddInboundTyped(ctx, xdto.InboundCreate(in))
}

func (d *sanaei) ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error) {
    return nil, nil
}

func (d *sanaei) CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
    orig, err := d.GetInboundTyped(ctx, inboundID)
    if err != nil { return xdto.Inbound{}, err }
    newPort := 0
    if opts.Port != nil { newPort = *opts.Port } else { newPort = randPort() }
    baseRemark := orig.Remark; if baseRemark == "" { baseRemark = "inb" }
    newRemark := baseRemark + "-copy-" + randSuffix(5)
    if opts.Remark != nil && *opts.Remark != "" { newRemark = *opts.Remark }
    // retry on conflict a couple of times
    for i:=0; i<2; i++ {
        inb, err := d.AddInboundTyped(ctx, xdto.InboundCreate{
            Remark: newRemark, Protocol: orig.Protocol, Port: newPort,
            Settings: orig.Settings, Sniffing: orig.Sniffing, StreamSettings: orig.StreamSettings,
        })
        if err == nil { return inb, nil }
        newPort = randPort(); newRemark = baseRemark + "-copy-" + randSuffix(6)
    }
    return d.AddInboundTyped(ctx, xdto.InboundCreate{
        Remark: newRemark, Protocol: orig.Protocol, Port: newPort,
        Settings: orig.Settings, Sniffing: orig.Sniffing, StreamSettings: orig.StreamSettings,
    })
}
