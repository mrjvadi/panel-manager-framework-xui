package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

type generic struct {
	sp   core.PanelSpec
	cli  *core.HTTP
	mu   sync.Mutex
	stat struct{ token string }
}

func ensureEndpoints(m map[string]string) map[string]string {
	if m == nil {
		m = map[string]string{}
	}
	for k, v := range defaultEndpoints {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}
	return m
}

func newGeneric(sp core.PanelSpec, opts ...core.Option) *generic {
	sp.Endpoints = ensureEndpoints(sp.Endpoints)
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
	return &generic{sp: sp, cli: cli}
}

func (d *generic) Name() string { return "xui.generic" }
func (d *generic) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	return "generic"
}
func (d *generic) Capabilities() core.Capabilities { return core.Capabilities{} }

func (d *generic) auth(req *http.Request) {
	if d.stat.token != "" {
		req.Header.Set("Authorization", "Bearer "+d.stat.token)
	}
}

func (d *generic) Login(ctx context.Context) error {
	body := map[string]any{"username": d.sp.Auth.Username, "password": d.sp.Auth.Password}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["login"], bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	var out map[string]any
	if err := core.DoJSON(ctx, d.cli, req, &out); err != nil {
		return err
	}
	if t, ok := out["token"].(string); ok {
		d.stat.token = t
	}
	return nil
}

func (d *generic) doJSON(ctx context.Context, req *http.Request, out any) error {
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
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return d.doJSON(ctx, req, out)
}
func (d *generic) putJSON(ctx context.Context, path string, payload any, out any) error {
	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewReader(b)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPut, d.cli.BaseURL+path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return d.doJSON(ctx, req, out)
}
func (d *generic) delete(ctx context.Context, path string) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
	return d.doJSON(ctx, req, nil)
}

func randPort() int { return 20000 + rand.Intn(40000) }
func randSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
