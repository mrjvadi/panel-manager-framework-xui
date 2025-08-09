package marzban

import (
    "bytes"
    "io"
    "sync"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    mdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/marzban"
)

const DriverName = "marzban"

func init() { core.Register(DriverName, New) }

type driver struct {
    mu   sync.Mutex
    sp   core.PanelSpec
    cli  *core.HTTP
    cap  core.Capabilities
    stat struct{ token string }
    runtimeVersion string
}
    mu   sync.Mutex

    sp   core.PanelSpec
    cli  *core.HTTP
    cap  core.Capabilities
    stat struct{ token string }
}

func New(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
    cap := core.Capabilities{
        UsersCRUD: true, InboundsCRUD: true, TrafficStats: true,
        UserSuspend: true, UserReset: true,
        Extra: map[core.Feature]bool{
            core.FeatureSubscriptions: true,
            core.FeatureUsersUsage:    true,
            core.FeatureSystemInfo:    true,
        },
    }
    cfg := core.collectDriverConfig(opts...)
	cli := core.NewHTTP(sp.BaseURL, sp.TLS.Insecure, chooseTimeout(cfg.Timeout, 30*time.Second), cfg.HTTPClient)
	if cfg.BreakerThresh > 0 { cli.Br = core.NewBreaker(cfg.BreakerThresh, cfg.BreakerCooldown) }
	if cfg.Retry.MaxAttempts > 0 { cli.Retry = cfg.Retry }
    // دیفالت‌ها
    def := map[string]string{
        "login":        "/api/admin/token",
        "listUsers":    "/api/admin/users",
        "listInbounds": "/api/admin/inbounds",
        // specials
        "user":                "/api/user",
        "userByName":          "/api/user/%s",
        "userReset":           "/api/user/%s/reset",
        "userRevoke":          "/api/user/%s/revoke_sub",
        "users":               "/api/users",
        "usersResetAll":       "/api/users/reset_all",
        "usersExpired":        "/api/users/expired",
        "userUsage":           "/api/user/%s/usage",
        "usersUsage":          "/api/users/usage",
        "subscription":        "/api/subscription",
        "subscriptionByID":    "/api/subscription/%s",
        "subscriptions":       "/api/subscriptions",
        "system":              "/api/system",
    }
    sp.Endpoints = core.MergeDefaults(def, sp.Endpoints)
    return &driver{ sp: sp, cli: cli, cap: cap }, nil
}

func (d *driver) Name() string                    { return DriverName }
func (d *driver) Version() string                 { if d.sp.Version != "" { return d.sp.Version }; if d.runtimeVersion != "" return d.runtimeVersion; return PluginVer }
func (d *driver) Capabilities() core.Capabilities { return d.cap }

// ===== Auth =====
func (d *driver) Login(ctx context.Context) error {
    form := url.Values{"grant_type": {"password"}}
    if b, ok := d.sp.Auth.(core.BasicAuth); ok { form.Set("username", b.Username); form.Set("password", b.Password) }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["login"], strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := d.cli.Client.Do(req); if err != nil { return err }
    defer resp.Body.Close()
    var body struct{ AccessToken string `json:"access_token"` }
    _ = json.NewDecoder(resp.Body).Decode(&body)
    if body.AccessToken != "" { d.stat.token = body.AccessToken }
    return nil
}

func (d *driver) auth(req *http.Request) { if d.stat.token != "" { req.Header.Set("Authorization", "Bearer "+d.stat.token) } }


// doJSON: ارسال درخواست با احراز هویت و هندل 401 (یک بار)
func (d *driver) doJSON(ctx context.Context, req *http.Request, out any) error {
    d.auth(req)
    if err := core.DoJSON(ctx, d.cli, req, out); err != nil {
        if core.IsHTTPStatus(err, http.StatusUnauthorized) {
            // re-login once
            d.mu.Lock()
            d.stat.token = ""
            _ = d.Login(ctx)
            d.mu.Unlock()
            // retry
            req2 := req.Clone(ctx)
            d.auth(req2)
            return core.DoJSON(ctx, d.cli, req2, out)
        }
        return err
    }
    return nil
}

func (d *driver) getJSON(ctx context.Context, path string, out any) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
    return d.doJSON(ctx, req, out)
}

func (d *driver) putJSON(ctx context.Context, path string, payload any, out any) error {
    var body io.Reader
    if payload != nil {
        b, _ := json.Marshal(payload)
        body = bytes.NewReader(b)
    }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPut, d.cli.BaseURL+path, body)
    if payload != nil { req.Header.Set("Content-Type", "application/json") }
    return d.doJSON(ctx, req, out)
}

func (d *driver) postJSON(ctx context.Context, path string, payload any, out any) error {
    var body io.Reader
    if payload != nil {
        b, _ := json.Marshal(payload)
        body = bytes.NewReader(b)
    }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, body)
    if payload != nil { req.Header.Set("Content-Type", "application/json") }
    return d.doJSON(ctx, req, out)
}

func (d *driver) delete(ctx context.Context, path string) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
    return d.doJSON(ctx, req, nil)
}

// ===== Core Lists =====
func (d *driver) ListUsers(ctx context.Context) ([]core.User, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["listUsers"], &body); err != nil { return nil, err }
    rows := extractArray(body, "users", "data", "items")
    out := make([]core.User, 0, len(rows))
    for _, r := range rows { out = append(out, mapUser(r)) }
    return out, nil
}

func (d *driver) ListInbounds(ctx context.Context) ([]core.Inbound, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["listInbounds"], &body); err != nil { return nil, err }
    rows := extractArray(body, "inbounds", "data", "items")
    out := make([]core.Inbound, 0, len(rows))
    for _, r := range rows { out = append(out, mapInbound(r)) }
    return out, nil
}

// ===== Users CRUD =====
func (d *driver) CreateUser(ctx context.Context, u core.User) (core.User, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    b, _ := json.Marshal(u)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["user"], bytes.NewReader(b))
    d.auth(req); req.Header.Set("Content-Type", "application/json")
    resp, err := d.cli.Client.Do(req); if err != nil { return u, err }
    defer resp.Body.Close()
    var out core.User; _ = json.NewDecoder(resp.Body).Decode(&out)
    if out.ID == "" { out = u }
    return out, nil
}

func (d *driver) UpdateUser(ctx context.Context, u core.User) (core.User, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["userByName"], u.Username)
    var out core.User
    if err := d.putJSON(ctx, path, u, &out); err != nil { return u, err }
    if out.ID == "" { out = u }
    return out, nil
}

func (d *driver) DeleteUser(ctx context.Context, id string) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["userByName"], id)
    return d.delete(ctx, path)
}

// ===== User Ops =====
func (d *driver) SuspendUser(ctx context.Context, id string) error { return core.ErrNotImplemented }
func (d *driver) ResumeUser(ctx context.Context, id string) error  { return core.ErrNotImplemented }
func (d *driver) ResetUserTraffic(ctx context.Context, id string) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["userReset"], id)
    return d.postJSON(ctx, path, nil, nil)
}

// ===== Extensions: Subscriptions =====
var _ ext.Subscription = (*driver)(nil)

func (d *driver) CreateSubscription(ctx context.Context, payload map[string]any) (map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var out map[string]any
    if err := d.postJSON(ctx, d.sp.Endpoints["subscription"], payload, &out); err != nil { return nil, err }
    return out, nil
}

func (d *driver) GetSubscription(ctx context.Context, id string) (map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["subscriptionByID"], id)
    var out map[string]any
    if err := d.getJSON(ctx, path, &out); err != nil { return nil, err }
    return out, nil
}

func (d *driver) ListSubscriptions(ctx context.Context) ([]map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["subscriptions"], &body); err != nil { return nil, err }
    return extractArray(body, "data", "items", "subscriptions"), nil
}

func (d *driver) DeleteSubscription(ctx context.Context, id string) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["subscriptionByID"], id)
    return d.delete(ctx, path)
}

func (d *driver) RevokeUserSubscription(ctx context.Context, username string) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["userRevoke"], username)
    return d.postJSON(ctx, path, nil, nil)
}

// ===== Extensions: Usage =====
var _ ext.Usage = (*driver)(nil)

func (d *driver) UserUsage(ctx context.Context, username string) (map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    path := fmt.Sprintf(d.sp.Endpoints["userUsage"], username)
    var out map[string]any
    if err := d.getJSON(ctx, path, &out); err != nil { return nil, err }
    return out, nil
}

func (d *driver) UsersUsage(ctx context.Context) ([]map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["usersUsage"], &body); err != nil { return nil, err }
    return extractArray(body, "data", "items", "users"), nil
}

func (d *driver) UsersExpired(ctx context.Context) ([]map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var body any
    if err := d.getJSON(ctx, d.sp.Endpoints["usersExpired"], &body); err != nil { return nil, err }
    return extractArray(body, "users", "data", "items"), nil
}

func (d *driver) ResetAllUsers(ctx context.Context) error {
    if d.stat.token == "" { _ = d.Login(ctx) }
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["usersResetAll"], nil)
    d.auth(req); resp, err := d.cli.Client.Do(req); if err != nil { return err }
    defer resp.Body.Close(); return nil
}

// ===== mapping helpers =====
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
func mapUser(r map[string]any) core.User {
    u := core.User{ ID: pickString(r, "id", "_id", "username"), Username: pickString(r, "username", "user", "name"), Raw: r }
    if v := pickInt64Ptr(r, "expire_at", "expire", "expiry_time"); v != nil { u.ExpireAt = v }
    if v := firstNonNilInt64Ptr(sumFields(r, []string{"up","down"}), pickInt64Ptr(r, "used_traffic", "traffic_used")); v != nil { u.TrafficUsed = v }
    if u.ID == "" { u.ID = u.Username }
    return u
}
func mapInbound(r map[string]any) core.Inbound {
    inb := core.Inbound{ ID: pickString(r, "id", "tag", "_id"), Type: pickString(r, "protocol", "type"), Remark: pickString(r, "remark", "tag"), Raw: r }
    if p := pickIntPtr(r, "port"); p != nil { inb.Port = p }
    return inb
}
func pickString(m map[string]any, keys ...string) string { for _, k := range keys { if v, ok := m[k]; ok { if s, ok := v.(string); ok { return s } } } ; return "" }
func pickInt64Ptr(m map[string]any, keys ...string) *int64 { for _, k := range keys { if v, ok := m[k]; ok { switch t := v.(type) { case float64: vv := int64(t); return &vv; case int64: return &t; case int: vv := int64(t); return &vv } } } ; return nil }
func pickIntPtr(m map[string]any, keys ...string) *int { for _, k := range keys { if v, ok := m[k]; ok { switch t := v.(type) { case float64: vv := int(t); return &vv; case int64: vv := int(t); return &vv; case int: return &t } } } ; return nil }
func sumFields(m map[string]any, keys []string) *int64 { var s int64; var f bool; for _, k := range keys { if v, ok := m[k]; ok { f = true; switch t := v.(type) { case float64: s += int64(t); case int64: s += t; case int: s += int64(t) } } } ; if !f { return nil }; return &s }
func firstNonNilInt64Ptr(a, b *int64) *int64 { if a != nil { return a }; return b }


var _ core.Connector = (*driver)(nil)

func (d *driver) Connect(ctx context.Context) error {
    // login once eager
    return d.Login(ctx)
}


func (d *driver) SystemInfo(ctx context.Context) (map[string]any, error) {
    if d.stat.token == "" { _ = d.Login(ctx) }
    var out map[string]any
    if err := d.getJSON(ctx, d.sp.Endpoints["system"], &out); err != nil { return nil, err }
    // cache version if present
    if v, ok := out["version"].(string); ok && v != "" {
        d.mu.Lock(); d.runtimeVersion = v; d.mu.Unlock()
    }
    return out, nil
}


func (d *driver) SystemInfoTyped(ctx context.Context) (mdto.SystemInfo, error) {
    m, err := d.SystemInfo(ctx)
    if err != nil { return mdto.SystemInfo{}, err }
    info := mdto.SystemInfo{ Raw: m }
    if v, ok := m["version"].(string); ok { info.Version = v }
    if u, ok := m["uptime"].(float64); ok { info.Uptime = int64(u) }
    return info, nil
}


func (d *driver) UsersUsageTyped(ctx context.Context) ([]mdto.UserUsage, error) {
    arr, err := d.UsersUsage(ctx)
    if err != nil { return nil, err }
    out := make([]mdto.UserUsage, 0, len(arr))
    for _, it := range arr {
        uu := mdto.UserUsage{}
        if s, ok := it["username"].(string); ok { uu.Username = s }
        if v, ok := it["up"].(float64); ok { uu.Up = int64(v) }
        if v, ok := it["down"].(float64); ok { uu.Down = int64(v) }
        uu.Total = uu.Up + uu.Down
        out = append(out, uu)
    }
    return out, nil
}

func (d *driver) UserUsageTyped(ctx context.Context, username string) (mdto.UserUsage, error) {
    m, err := d.UserUsage(ctx, username)
    if err != nil { return mdto.UserUsage{}, err }
    uu := mdto.UserUsage{ Username: username }
    if v, ok := m["up"].(float64); ok { uu.Up = int64(v) }
    if v, ok := m["down"].(float64); ok { uu.Down = int64(v) }
    uu.Total = uu.Up + uu.Down
    return uu, nil
}

func (d *driver) ListSubscriptionsTyped(ctx context.Context) ([]mdto.Subscription, error) {
    arr, err := d.ListSubscriptions(ctx)
    if err != nil { return nil, err }
    out := make([]mdto.Subscription, 0, len(arr))
    for _, it := range arr {
        s := mdto.Subscription{ Raw: it }
        if v, ok := it["id"].(string); ok { s.ID = v }
        if v, ok := it["username"].(string); ok { s.Username = v }
        if v, ok := it["link"].(string); ok { s.Link = v }
        out = append(out, s)
    }
    return out, nil
}


type healthOK struct{}

var _ core.Connector = (*driver)(nil)
var _ core.HealthChecker = (*driver)(nil)
var _ ext.MarzbanTyped = (*driver)(nil)

func (d *driver) Connect(ctx context.Context) error {
    if err := d.Login(ctx); err != nil { return err }
    // discover version
    _, _ = d.SystemInfoTyped(ctx)
    return nil
}

func (d *driver) Health(ctx context.Context) error {
    // simple system call
    _, err := d.SystemInfo(ctx)
    return err
}
