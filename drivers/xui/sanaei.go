package xui

import (
    "bytes"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const SanaeiName = "xui.sanaei"

func init() { core.Register(SanaeiName, NewSanaei) }

type sanaei struct { *generic }

func NewSanaei(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    if sp.Endpoints == nil { sp.Endpoints = map[string]string{} }
    def := map[string]string{
        "login":            "/login",
        "listUsers":        "/xui/user/list",
        "listInbounds":     "/panel/api/inbounds/list",
        "getInbound":       "/panel/api/inbounds/get/%d",
        "addInbound":       "/panel/api/inbounds/add",
        "updateInbound":    "/panel/api/inbounds/update/%d",
        "deleteInbound":    "/panel/api/inbounds/del/%d",
        "addClient":        "/panel/api/inbounds/addClient",
        "deleteClient":     "/panel/api/inbounds/delClient/%d",
        "clientTraffEmail": "/panel/api/inbounds/getClientTraffics/%s",
        "clientTraffID":    "/panel/api/inbounds/getClientTrafficsById/%s",
        "resetAllTraffic":  "/panel/api/inbounds/resetAllTraffic",
        "clientIPs":        "/panel/api/inbounds/clientIps/%s",
    }
    sp.Endpoints = core.MergeDefaults(def, sp.Endpoints)
    if sp.Version == "" { sp.Version = Version3XUI }
    g := newGeneric(sp, opts...)
    return &sanaei{ g }, nil
}

func (d *sanaei) Name() string    { return SanaeiName }
func (d *sanaei) Version() string { if d.sp.Version != "" { return d.sp.Version }; return Version3XUI }

// === افزونه‌های مخصوص X-UI: Inbounds/Clients ===
var _ ext.InboundsAdmin = (*sanaei)(nil)

func (d *generic) GetInbound(ctx context.Context, inboundID int) (map[string]any, error) {
    path := fmt.Sprintf(d.sp.Endpoints["getInbound"], inboundID)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) AddInbound(ctx context.Context, payload map[string]any) (map[string]any, error) {
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["addInbound"], bytes.NewReader(b))
    d.auth(req); req.Header.Set("Content-Type", "application/json")
    resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) UpdateInbound(ctx context.Context, inboundID int, payload map[string]any) (map[string]any, error) {
    b, _ := json.Marshal(payload); path := fmt.Sprintf(d.sp.Endpoints["updateInbound"], inboundID)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, bytes.NewReader(b))
    d.auth(req); req.Header.Set("Content-Type", "application/json")
    resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) DeleteInbound(ctx context.Context, inboundID int) error {
    path := fmt.Sprintf(d.sp.Endpoints["deleteInbound"], inboundID)
    req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return err }
    defer resp.Body.Close(); return nil
}

func (d *generic) AddClient(ctx context.Context, payload map[string]any) (map[string]any, error) {
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["addClient"], bytes.NewReader(b))
    d.auth(req); req.Header.Set("Content-Type", "application/json")
    resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) DeleteClient(ctx context.Context, clientID int) error {
    path := fmt.Sprintf(d.sp.Endpoints["deleteClient"], clientID)
    req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return err }
    defer resp.Body.Close(); return nil
}

func (d *generic) ClientTrafficByEmail(ctx context.Context, email string) (map[string]any, error) {
    path := fmt.Sprintf(d.sp.Endpoints["clientTraffEmail"], email)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) ClientTrafficByID(ctx context.Context, uuid string) (map[string]any, error) {
    path := fmt.Sprintf(d.sp.Endpoints["clientTraffID"], uuid)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var out map[string]any; _ = json.NewDecoder(resp.Body).Decode(&out); return out, nil
}

func (d *generic) ResetAllTraffic(ctx context.Context) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["resetAllTraffic"], nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return err }
    defer resp.Body.Close(); return nil
}

func (d *generic) ClientIPs(ctx context.Context, email string) ([]string, error) {
    path := fmt.Sprintf(d.sp.Endpoints["clientIPs"], email)
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return nil, err }
    defer resp.Body.Close(); var body any; _ = json.NewDecoder(resp.Body).Decode(&body)
    arr := extractArray(body, "ips", "data", "items"); out := make([]string, 0, len(arr))
    for _, it := range arr { if s, ok := it["ip"].(string); ok { out = append(out, s) } }
    return out, nil
}


// Typed API for X-UI (Sanaei)
var _ ext.XUITyped = (*sanaei)(nil)

func (d *sanaei) GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error) {
    m, err := d.GetInbound(ctx, inboundID)
    if err != nil { return xdto.Inbound{}, err }
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
    m, err := d.AddInbound(ctx, map[string]any{
        "remark": in.Remark, "protocol": in.Protocol, "port": in.Port,
        "settings": in.Settings, "sniffing": in.Sniffing, "streamSettings": in.StreamSettings,
    })
    if err != nil { return xdto.Inbound{}, err }
    return d.GetInboundTyped(ctx, int(asFloat(m["id"])))
}

func (d *sanaei) UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error) {
    _, err := d.UpdateInbound(ctx, inboundID, map[string]any{
        "remark": in.Remark, "protocol": in.Protocol, "port": in.Port,
        "settings": in.Settings, "sniffing": in.Sniffing, "streamSettings": in.StreamSettings,
    })
    if err != nil { return xdto.Inbound{}, err }
    return d.GetInboundTyped(ctx, inboundID)
}

func (d *sanaei) ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error) {
    m, err := d.ClientTrafficByEmail(ctx, email)
    if err != nil { return nil, err }
    // معمولاً خروجی می‌تواند آرایه‌ای داخل فیلدی باشد؛ از generic استخراج می‌کنیم
    var arr []map[string]any
    if a, ok := m["items"].([]any); ok {
        for _, it := range a {
            if mm, ok := it.(map[string]any); ok { arr = append(arr, mm) }
        }
    } else if a, ok := m["data"].([]any); ok {
        for _, it := range a {
            if mm, ok := it.(map[string]any); ok { arr = append(arr, mm) }
        }
    } else if mm, ok := m["traffic"].([]any); ok {
        for _, it := range mm {
            if m2, ok := it.(map[string]any); ok { arr = append(arr, m2) }
        }
    }
    out := make([]xdto.TrafficRecord, 0, len(arr))
    for _, it := range arr {
        tr := xdto.TrafficRecord{}
        if s, ok := it["email"].(string); ok { tr.Email = s }
        tr.Up = int64(asFloat(it["up"])); tr.Down = int64(asFloat(it["down"])); tr.Total = tr.Up + tr.Down
        out = append(out, tr)
    }
    return out, nil
}

// helper
func asFloat(v any) float64 { if f, ok := v.(float64); ok { return f }; if i, ok := v.(int); ok { return float64(i) }; if i, ok := v.(int64); ok { return float64(i) }; return 0 }


// CloneInboundTyped: کلون‌کردن یک inbound با امکان override پورت/اسم
func (d *sanaei) CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
    // 1) دریافت inbound اصلی
    orig, err := d.GetInboundTyped(ctx, inboundID)
    if err != nil { return xdto.Inbound{}, err }

    // 2) آماده‌سازی مقادیر جدید
    newPort := 0
    if opts.Port != nil { newPort = *opts.Port } else { newPort = randPort() }

    baseRemark := orig.Remark
    if baseRemark == "" { baseRemark = "inb" }
    newRemark := baseRemark + "-copy-" + randSuffix(5)
    if opts.Remark != nil && *opts.Remark != "" { newRemark = *opts.Remark }

    // 3) ساخت payload
    payload := xdto.InboundCreate{
        Remark: newRemark,
        Protocol: orig.Protocol,
        Port: newPort,
        Settings: orig.Settings,
        Sniffing: orig.Sniffing,
        StreamSettings: orig.StreamSettings,
    }

    // 4) ارسال - با تلاش مجدد محدود روی برخورد (مثل 409)
    attempts := 3
    for i := 0; i < attempts; i++ {
        inb, err := d.AddInboundTyped(ctx, payload)
        if err == nil {
            return inb, nil
        }
        // اگر خطای HTTP با کد تداخل بود، اسم/پورت را عوض کن و دوباره امتحان کن
        if he, ok := err.(*core.HTTPError); ok && (he.Code == 400 || he.Code == 409) {
            payload.Port = randPort()
            payload.Remark = baseRemark + "-copy-" + randSuffix(6)
            continue
        }
        return xdto.Inbound{}, err
    }
    // آخرین تلاش
    return d.AddInboundTyped(ctx, payload)
}
