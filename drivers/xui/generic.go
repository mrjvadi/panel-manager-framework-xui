package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

// داخل struct:
type generic struct {
	sp         core.PanelSpec
	cli        *core.HTTP
	mu         sync.Mutex
	cookieName string
	cookieVal  string
	stat       struct {
		cookie string // <- به جای token
	}
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

// کوکی را روی درخواست بفرست:
func (d *generic) auth(req *http.Request) {
	if d.stat.cookie != "" {
		req.Header.Set("Cookie", d.stat.cookie)
	}
	if d.cookieName != "" && d.cookieVal != "" {
		req.Header.Set("Cookie", d.cookieName+"="+d.cookieVal)
	}
}

// لاگین با فرم‌داده (x-www-form-urlencoded) + گرفتن کوکی سشن
func (d *generic) Login(ctx context.Context) error {
	// form-urlencoded: username=...&password=...&twoFactorCode=
	form := url.Values{}
	form.Set("username", d.sp.Auth.Username)
	form.Set("password", d.sp.Auth.Password)
	form.Set("twoFactorCode", "")

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		d.cli.BaseURL+d.sp.Endpoints["login"], strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &core.HTTPError{Code: resp.StatusCode, Body: string(b)}
	}

	name, val := d.pickSessionCookie(resp)
	if name == "" || val == "" {
		return fmt.Errorf("session cookie not found (tried 3x-ui/x-ui/xui/*session*)")
	}

	d.cookieName, d.cookieVal = name, val
	return nil
}

func looksLikeSession(name string) bool {
	n := strings.ToLower(name)
	return n == "3x-ui" || n == "x-ui" || n == "xui" || strings.Contains(n, "session")
}

func (d *generic) pickSessionCookie(resp *http.Response) (name, val string) {
	// 1) از خود پاسخ
	for _, ck := range resp.Cookies() {
		if looksLikeSession(ck.Name) {
			return ck.Name, ck.Value
		}
	}
	// 2) از کوکی‌جار روی URL واقعی درخواست (نه BaseURL)
	if d.cli.Client.Jar != nil && resp != nil && resp.Request != nil && resp.Request.URL != nil {
		for _, ck := range d.cli.Client.Jar.Cookies(resp.Request.URL) {
			if looksLikeSession(ck.Name) {
				return ck.Name, ck.Value
			}
		}
	}
	// 3) از Set-Cookie خام (fallback)
	for _, raw := range resp.Header.Values("Set-Cookie") {
		parts := strings.SplitN(raw, ";", 2)
		if len(parts) > 0 {
			kv := strings.SplitN(parts[0], "=", 2)
			if len(kv) == 2 && looksLikeSession(kv[0]) {
				return kv[0], kv[1]
			}
		}
	}
	return "", ""
}

// رپر روی DoJSON که اگر JSON نبود/به صفحه لاگین هدایت شد، یک‌بار re-login کند
func (d *generic) doJSON(ctx context.Context, req *http.Request, out any) error {
	if d.stat.cookie == "" {
		_ = d.Login(ctx) // تلاش اولیه
	}
	d.auth(req)
	if err := core.DoJSON(ctx, d.cli, req, out); err != nil {
		// اگر Unauthorized بود، re-login و تکرار
		if core.IsHTTPStatus(err, http.StatusUnauthorized) {
			d.mu.Lock()
			d.stat.cookie = ""
			_ = d.Login(ctx)
			d.mu.Unlock()
			req2 := req.Clone(ctx)
			d.auth(req2)
			return core.DoJSON(ctx, d.cli, req2, out)
		}
		// اگر decode JSON شکست خورد و به احتمال زیاد HTML لاگین برگشته
		if strings.Contains(err.Error(), "decode json") || strings.Contains(err.Error(), "invalid character '<'") {
			d.mu.Lock()
			d.stat.cookie = ""
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
