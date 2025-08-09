package xui

import (
    "context"
    "fmt"
    "net/http"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const SanaeiName = "xui.sanaei"
const PluginVer = "v2.6.1"

func init() { core.Register(SanaeiName, NewSanaei) }

type sanaei struct{ *generic }

func NewSanaei(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    g := newGeneric(sp, opts...)
    return &sanaei{ g }, nil
}

func (d *sanaei) Name() string    { return SanaeiName }
func (d *sanaei) Version() string { if d.sp.Version != "" { return d.sp.Version }; return PluginVer }
func (d *sanaei) Capabilities() core.Capabilities { return core.Capabilities{} }

func (d *sanaei) ListUsers(ctx context.Context) ([]core.User, error) { return nil, nil }

func (d *sanaei) ListInbounds(ctx context.Context) ([]core.Inbound, error) {
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["listInbounds"], &body); err != nil { return nil, err }
    return nil, nil
}

// Typed API
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
    path := fmt.Sprintf(d.sp.Endpoints["updateInbound"], inboundID)
    var m map[string]any
    if err := d.postJSON(ctx, path, map[string]any{
        "remark": in.Remark, "protocol": in.Protocol, "port": in.Port,
        "settings": in.Settings, "sniffing": in.Sniffing, "streamSettings": in.StreamSettings,
    }, &m); err != nil {
        if core.IsHTTPStatus(err, http.StatusNotFound) {
            return xdto.Inbound{}, err
        }
        return xdto.Inbound{}, err
    }
    if v, ok := m["id"].(float64); ok && int(v) > 0 { return d.GetInboundTyped(ctx, int(v)) }
    return d.GetInboundTyped(ctx, inboundID)
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
    return d.AddInboundTyped(ctx, xdto.InboundCreate{
        Remark: newRemark, Protocol: orig.Protocol, Port: newPort,
        Settings: orig.Settings, Sniffing: orig.Sniffing, StreamSettings: orig.StreamSettings,
    })
}
