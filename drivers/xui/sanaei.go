package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

const SanaeiName = "xui.sanaei"

func init() { core.Register(SanaeiName, NewSanaei) }

type sanaei struct{ generic }

func NewSanaei(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
	if sp.Endpoints == nil {
		sp.Endpoints = map[string]string{}
	}
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
	if sp.Version == "" {
		sp.Version = Version3XUI
	}
	g := newGeneric(sp, opts...)
	return &sanaei{*g}, nil
}

func (d *sanaei) Name() string { return SanaeiName }
func (d *sanaei) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	return Version3XUI
}

// === افزونه‌های مخصوص X-UI: Inbounds/Clients ===
var _ ext.InboundsAdmin = (*sanaei)(nil)

func (d *generic) GetInbound(ctx context.Context, inboundID int) (map[string]any, error) {
	path := fmt.Sprintf(d.sp.Endpoints["getInbound"], inboundID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) AddInbound(ctx context.Context, payload map[string]any) (map[string]any, error) {
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["addInbound"], bytes.NewReader(b))
	d.auth(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) UpdateInbound(ctx context.Context, inboundID int, payload map[string]any) (map[string]any, error) {
	b, _ := json.Marshal(payload)
	path := fmt.Sprintf(d.sp.Endpoints["updateInbound"], inboundID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, bytes.NewReader(b))
	d.auth(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) DeleteInbound(ctx context.Context, inboundID int) error {
	path := fmt.Sprintf(d.sp.Endpoints["deleteInbound"], inboundID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (d *generic) AddClient(ctx context.Context, payload map[string]any) (map[string]any, error) {
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["addClient"], bytes.NewReader(b))
	d.auth(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) DeleteClient(ctx context.Context, clientID int) error {
	path := fmt.Sprintf(d.sp.Endpoints["deleteClient"], clientID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (d *generic) ClientTrafficByEmail(ctx context.Context, email string) (map[string]any, error) {
	path := fmt.Sprintf(d.sp.Endpoints["clientTraffEmail"], email)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) ClientTrafficByID(ctx context.Context, uuid string) (map[string]any, error) {
	path := fmt.Sprintf(d.sp.Endpoints["clientTraffID"], uuid)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}

func (d *generic) ResetAllTraffic(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.sp.Endpoints["resetAllTraffic"], nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (d *generic) ClientIPs(ctx context.Context, email string) ([]string, error) {
	path := fmt.Sprintf(d.sp.Endpoints["clientIPs"], email)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, d.cli.BaseURL+path, nil)
	d.auth(req)
	resp, err := d.cli.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var body any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	arr := extractArray(body, "ips", "data", "items")
	out := make([]string, 0, len(arr))
	for _, it := range arr {
		if s, ok := it["ip"].(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}
