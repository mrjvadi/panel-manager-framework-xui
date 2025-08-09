package xui

import (
    "bytes"
    "io"
    "sync"
    "math/rand"
    "time"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
)

const GenericName = "xui.generic"

func init() { core.Register(GenericName, NewGeneric) }

type generic struct {
    mu   sync.Mutex
    sp   core.PanelSpec
    cli  *core.HTTP
    cap  core.Capabilities
    stat struct{ token string }
}

func newGeneric(sp core.PanelSpec, opts ...core.Option) *generic {
    cfg := core.collectDriverConfig(opts...)
	cli := core.NewHTTP(sp.BaseURL, sp.TLS.Insecure, chooseTimeout(cfg.Timeout, 30*time.Second), cfg.HTTPClient)
	if cfg.BreakerThresh > 0 { cli.Br = core.NewBreaker(cfg.BreakerThresh, cfg.BreakerCooldown) }
	if cfg.Retry.MaxAttempts > 0 { cli.Retry = cfg.Retry }
    cap := core.Capabilities{ UsersCRUD: true, InboundsCRUD: true, Extra: map[core.Feature]bool{ core.FeatureXUIClients: true } }
    return &generic{ sp: sp, cli: cli, cap: cap }
}

func NewGeneric(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    return newGeneric(sp, opts...), nil
}

func (d *generic) Name() string                    { return GenericName }
func (d *generic) Version() string                 { if d.sp.Version != "" { return d.sp.Version }; return "generic" }
func (d *generic) Capabilities() core.Capabilities { return d.cap }

func (d *generic) Login(ctx context.Context) error {
    path := d.sp.Endpoints["login"]
    if path == "" { path = "/login" }
    body := map[string]string{}
    if b, ok := d.sp.Auth.(core.BasicAuth); ok { body["username"], body["password"] = b.Username, b.Password }
    bts, _ := json.Marshal(body)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, bytes.NewReader(bts))
    req.Header.Set("Content-Type", "application/json")
    resp, err := d.cli.Client.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    var res any; _ = json.NewDecoder(resp.Body).Decode(&res)
    if m, ok := res.(map[string]any); ok {
        if t, ok := m["accessToken"].(string); ok && t != "" { d.stat.token = t }
        if t, ok := m["token"].(string); ok && t != "" { d.stat.token = t }
    }
    return nil
}

func (d *generic) auth(req *http.Request) { if d.stat.token != "" { req.Header.Set("Authorization", "Bearer "+d.stat.token) } }

func (d *generic) ListUsers(ctx context.Context) ([]core.User, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := d.sp.Endpoints["listUsers"]
    if path == "" { path = "/xui/user/list" }
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req)
    resp, err := d.cli.Client.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    b, _ := io.ReadAll(resp.Body)
    var body any; _ = json.Unmarshal(b, &body)
    rows := extractArray(body, "data", "users", "items")
    out := make([]core.User, 0, len(rows))
    for _, r := range rows {
        u := core.User{ ID: pickString(r, "id", "_id", "username"), Username: pickString(r, "username", "user"), Raw: r }
        if v := pickInt64Ptr(r, "expire", "expire_at"); v != nil { u.ExpireAt = v }
        if v := firstNonNilInt64Ptr(sumFields(r, []string{"up","down"}), pickInt64Ptr(r, "traffic_used")); v != nil { u.TrafficUsed = v }
        if u.ID == "" { u.ID = u.Username }
        out = append(out, u)
    }
    return out, nil
}

func (d *generic) ListInbounds(ctx context.Context) ([]core.Inbound, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := d.sp.Endpoints["listInbounds"]
    if path == "" { path = "/xui/inbound/list" }
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    d.auth(req)
    resp, err := d.cli.Client.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    var body any
    if err := json.NewDecoder(resp.Body).Decode(&body); err != nil { return nil, err }
    rows := extractArray(body, "data", "items")
    out := make([]core.Inbound, 0, len(rows))
    for _, r := range rows {
        inb := core.Inbound{ ID: pickString(r, "id", "_id"), Type: pickString(r, "protocol", "type"), Remark: pickString(r, "remark", "tag"), Raw: r }
        if p := pickIntPtr(r, "port"); p != nil { inb.Port = p }
        out = append(out, inb)
    }
    return out, nil
}

func (d *generic) CreateUser(ctx context.Context, u core.User) (core.User, error)      { return u, core.ErrNotImplemented }
func (d *generic) UpdateUser(ctx context.Context, u core.User) (core.User, error)      { return u, core.ErrNotImplemented }
func (d *generic) DeleteUser(ctx context.Context, id string) error                    { return core.ErrNotImplemented }
func (d *generic) SuspendUser(ctx context.Context, id string) error                   { return core.ErrNotImplemented }
func (d *generic) ResumeUser(ctx context.Context, id string) error                    { return core.ErrNotImplemented }
func (d *generic) ResetUserTraffic(ctx context.Context, id string) error              { return core.ErrNotImplemented }
func (d *generic) CreateInbound(ctx context.Context, in core.Inbound) (core.Inbound, error) { return in, core.ErrNotImplemented }
func (d *generic) UpdateInbound(ctx context.Context, in core.Inbound) (core.Inbound, error) { return in, core.ErrNotImplemented }
func (d *generic) DeleteInbound(ctx context.Context, id string) error                  { return core.ErrNotImplemented }

// helpers مشترک
func extractArray(body any, keys ...string) []map[string]any {
    if m, ok := body.(map[string]any); ok {
        for _, k := range keys {
            if v, ok := m[k]; ok {
                if arr, ok := v.([]any); ok {
                    out := make([]map[string]any, 0, len(arr))
                    for _, it := range arr {
                        if mm, ok := it.(map[string]any); ok { out = append(out, mm) }
                    }
                    return out
                }
            }
        }
    }
    if arr, ok := body.([]any); ok {
        out := make([]map[string]any, 0, len(arr))
        for _, it := range arr {
            if mm, ok := it.(map[string]any); ok { out = append(out, mm) }
        }
        return out
    }
    return nil
}
func pickString(m map[string]any, keys ...string) string { for _, k := range keys { if v, ok := m[k]; ok { if s, ok := v.(string); ok { return s } } } ; return "" }
func pickInt64Ptr(m map[string]any, keys ...string) *int64 { for _, k := range keys { if v, ok := m[k]; ok { switch t := v.(type) { case float64: vv := int64(t); return &vv; case int64: return &t; case int: vv := int64(t); return &vv } } } ; return nil }
func pickIntPtr(m map[string]any, keys ...string) *int { for _, k := range keys { if v, ok := m[k]; ok { switch t := v.(type) { case float64: vv := int(t); return &vv; case int64: vv := int(t); return &vv; case int: return &t } } } ; return nil }
func sumFields(m map[string]any, keys []string) *int64 { var s int64; var f bool; for _, k := range keys { if v, ok := m[k]; ok { f = true; switch t := v.(type) { case float64: s += int64(t); case int64: s += t; case int: s += int64(t) } } } ; if !f { return nil }; return &s }
func firstNonNilInt64Ptr(a, b *int64) *int64 { if a != nil { return a }; return b }


func (d *generic) doJSON(ctx context.Context, req *http.Request, out any) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    // attach auth and call core.DoJSON with retry/breaker
    d.auth(req)
    if err := core.DoJSON(ctx, d.cli, req, out); err != nil {
        if core.IsHTTPStatus(err, http.StatusUnauthorized) {
            d.mu.Lock(); d.stat.token = ""; _ = d.Login(ctx); d.mu.Unlock()
            req2 := req.Clone(ctx); d.auth(req2)
            return core.DoJSON(ctx, d.cli, req2, out)
        }
        return err
    }
    return nil
}
func (d *generic) getJSON(ctx context.Context, path string, out any) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    return d.doJSON(ctx, req, out)
}
func (d *generic) postJSON(ctx context.Context, path string, payload any, out any) error {
    var body io.Reader
    if payload != nil {
        b, _ := json.Marshal(payload)
        body = bytes.NewReader(b)
    }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, body)
    if payload != nil { req.Header.Set("Content-Type", "application/json") }
    return d.doJSON(ctx, req, out)
}
func (d *generic) putJSON(ctx context.Context, path string, payload any, out any) error {
    var body io.Reader
    if payload != nil {
        b, _ := json.Marshal(payload)
        body = bytes.NewReader(b)
    }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPut, d.cli.BaseURL+path, body)
    if payload != nil { req.Header.Set("Content-Type", "application/json") }
    return d.doJSON(ctx, req, out)
}
func (d *generic) delete(ctx context.Context, path string) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
    return d.doJSON(ctx, req, nil)
}


var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func randPort() int { return 20000 + rng.Intn(40000) }

func randSuffix(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, n)
    for i := range b { b[i] = letters[rng.Intn(len(letters))] }
    return string(b)
}
