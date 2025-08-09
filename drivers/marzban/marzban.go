package marzban

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

const DriverName = "marzban"
const PluginVer = "v0.8.4"

func init() { core.Register(DriverName, New) }

type driver struct {
	sp             core.PanelSpec
	cli            *core.HTTP
	mu             sync.Mutex
	stat           struct{ token string }
	runtimeVersion string
}

func New(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
	if sp.Endpoints == nil {
		sp.Endpoints = map[string]string{}
	}
	if sp.Endpoints["login"] == "" {
		sp.Endpoints["login"] = "/api/admin/token"
	}
	if sp.Endpoints["listUsers"] == "" {
		sp.Endpoints["listUsers"] = "/api/admin/users"
	}
	if sp.Endpoints["system"] == "" {
		sp.Endpoints["system"] = "/api/system"
	}
	cfg := core.CollectDriverConfig(opts...)
	cli := core.NewHTTP(sp.BaseURL, sp.TLS.Insecure, core.ChooseTimeout(cfg.Timeout, 30*time.Second), cfg.HTTPClient)
	if cfg.BreakerThresh > 0 {
		cli.Br = core.NewBreaker(cfg.BreakerThresh, cfg.BreakerCooldown)
	}
	if cfg.Retry.MaxAttempts > 0 {
		cli.Retry = cfg.Retry
	}
	if cfg.Logger != nil {
		cli.Log = cfg.Logger
	}
	return &driver{sp: sp, cli: cli}, nil
}

func (d *driver) Name() string { return DriverName }
func (d *driver) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	if d.runtimeVersion != "" {
		return d.runtimeVersion
	}
	return PluginVer
}
func (d *driver) Capabilities() core.Capabilities { return core.Capabilities{} }

func (d *driver) auth(req *http.Request) {
	if d.stat.token != "" {
		req.Header.Set("Authorization", "Bearer "+d.stat.token)
	}
}

func (d *driver) Login(ctx context.Context) error {
	body := map[string]any{"username": d.sp.Auth.Username, "password": d.sp.Auth.Password}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["login"], bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	var out map[string]any
	if err := core.DoJSON(ctx, d.cli, req, &out); err != nil {
		return err
	}
	if t, ok := out["access_token"].(string); ok {
		d.stat.token = t
	}
	return nil
}

func (d *driver) doJSON(ctx context.Context, req *http.Request, out any) error {
	if d.stat.token == "" {
		_ = d.Login(ctx)
	}
	d.auth(req)
	if err := core.DoJSON(ctx, d.cli, req, out); err != nil {
		if core.IsHTTPStatus(err, http.StatusUnauthorized) {
			d.mu.Lock()
			d.stat.token = ""
			_ = d.Login(ctx)
			d.mu.Unlock()
			req2 := req.Clone(ctx)
			d.auth(req2)
			return core.DoJSON(ctx, d.cli, req2, out)
		}
		return err
	}
	return nil
}

func (d *driver) ListUsers(ctx context.Context) ([]core.User, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+d.sp.Endpoints["listUsers"], nil)
	var body map[string]any
	if err := d.doJSON(ctx, req, &body); err != nil {
		return nil, err
	}
	var out []core.User
	if arr, ok := body["users"].([]any); ok {
		for _, it := range arr {
			if m, ok := it.(map[string]any); ok {
				u := core.User{}
				if s, ok := m["username"].(string); ok {
					u.Username = s
				}
				if v, ok := m["up"].(float64); ok {
					u.Up = int64(v)
				}
				if v, ok := m["down"].(float64); ok {
					u.Down = int64(v)
				}
				out = append(out, u)
			}
		}
	}
	return out, nil
}
func (d *driver) ListInbounds(ctx context.Context) ([]core.Inbound, error) { return nil, nil }

func (d *driver) SystemInfo(ctx context.Context) (map[string]any, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+d.sp.Endpoints["system"], nil)
	var out map[string]any
	if err := d.doJSON(ctx, req, &out); err != nil {
		return nil, err
	}
	if v, ok := out["version"].(string); ok {
		d.runtimeVersion = v
	}
	return out, nil
}
