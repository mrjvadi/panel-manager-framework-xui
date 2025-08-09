package xui

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

// ensure the driver exposes this extension
var _ ext.XUIShallowClone = (*sanaei)(nil)

// ---- helpers (local and safe) ----

func (d *sanaei) ep(key, def string) string {
	if s := d.sp.Endpoints[key]; s != "" {
		return s
	}
	return def
}

func anyToInt(a any) int {
	switch v := a.(type) {
	case float64:
		return int(v)
	case json.Number:
		i, _ := strconv.Atoi(string(v))
		return i
	case string:
		i, _ := strconv.Atoi(v)
		return i
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

// Extracts the main object from various response structures
func takeObj(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	if obj, ok := m["obj"].(map[string]any); ok {
		return obj
	}
	if arr, ok := m["obj"].([]any); ok && len(arr) > 0 {
		if mm, ok2 := arr[len(arr)-1].(map[string]any); ok2 {
			return mm
		}
	}
	return m
}

func (d *sanaei) getInboundRaw(ctx context.Context, id int) (map[string]any, error) {
	path := d.ep("getInbound", "/panel/api/inbounds/get/%d")
	url := d.cli.BaseURL + fmt.Sprintf(path, id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	var resp map[string]any
	if err := d.doJSON(ctx, req, &resp); err != nil {
		return nil, err
	}
	return takeObj(resp), nil
}

func (d *sanaei) listInboundsRaw(ctx context.Context) ([]map[string]any, error) {
	url := d.cli.BaseURL + d.ep("listInbounds", "/panel/api/inbounds/list")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	var resp any
	if err := d.doJSON(ctx, req, &resp); err != nil {
		return nil, err
	}

	out := []map[string]any{}
	switch v := resp.(type) {
	case map[string]any:
		if obj, ok := v["obj"]; ok {
			if arr, ok := obj.([]any); ok {
				for _, e := range arr {
					if mm, ok2 := e.(map[string]any); ok2 {
						out = append(out, mm)
					}
				}
			}
		}
	case []any:
		for _, e := range v {
			if mm, ok2 := e.(map[string]any); ok2 {
				out = append(out, mm)
			}
		}
	}
	return out, nil
}

func (d *sanaei) resolveNewInbound(ctx context.Context, beforeIDs map[int]struct{}, remark string, port int) (xdto.Inbound, bool) {
	after, err := d.listInboundsRaw(ctx)
	if err != nil {
		return xdto.Inbound{}, false
	}
	best := xdto.Inbound{}
	bestID := 0
	for _, m := range after {
		id := anyToInt(m["id"])
		if _, existed := beforeIDs[id]; existed {
			continue
		}
		rm, _ := m["remark"].(string)
		prt := anyToInt(m["port"])
		if remark != "" && strings.EqualFold(rm, remark) {
			return xdto.Inbound{ID: id, Remark: rm, Port: prt, Raw: m}, true
		}
		if port > 0 && prt == port {
			return xdto.Inbound{ID: id, Remark: rm, Port: prt, Raw: m}, true
		}
		if id > bestID {
			bestID = id
			best = xdto.Inbound{ID: id, Remark: rm, Port: prt, Raw: m}
		}
	}
	if bestID > 0 {
		return best, true
	}
	return xdto.Inbound{}, false
}

// Generates a new UUID v4.
func newUUID() (string, error) {
	b := make([]byte, 16)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func randSlug(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	for i := 0; i < n; i++ {
		ri, _ := crand.Int(crand.Reader, big.NewInt(int64(len(letters))))
		b.WriteByte(letters[ri.Int64()])
	}
	return b.String()
}


// ---- Main implementation: CloneInboundShallow ----
func (d *sanaei) CloneInboundShallow(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
	orig, err := d.getInboundRaw(ctx, inboundID)
	if err != nil {
		return xdto.Inbound{}, err
	}

	// Remark
	var newRemark string
	if opts.Remark != nil && *opts.Remark != "" {
		newRemark = *opts.Remark
	} else {
		baseRemark, _ := orig["remark"].(string)
		if baseRemark == "" {
			baseRemark = "inb"
		}
		newRemark = baseRemark + "-" + randSlug(6)
	}

	// Port
	var port int
	if opts.Port != nil {
		port = *opts.Port
	} else {
		origPort := anyToInt(orig["port"])
		port = 20000 + mrand.Intn(40000)
		if port == origPort {
			port = 20000 + mrand.Intn(40000)
		}
	}

	// Get existing IDs to resolve the new one later
	before, _ := d.listInboundsRaw(ctx)
	beforeIDs := make(map[int]struct{}, len(before))
	for _, m := range before {
		if id := anyToInt(m["id"]); id > 0 {
			beforeIDs[id] = struct{}{}
		}
	}

	// Build the form with the correct client structure inside settings
	buildForm := func(rem string, p int) (url.Values, error) {
		v := url.Values{}
		v.Set("up", "0")
		v.Set("down", "0")
		v.Set("total", "0")
		v.Set("remark", rem)
		v.Set("enable", "true")
		v.Set("expiryTime", "0")
		v.Set("listen", "")
		v.Set("port", strconv.Itoa(p))
		if proto, _ := orig["protocol"].(string); proto != "" {
			v.Set("protocol", proto)
		} else {
			v.Set("protocol", "vless")
		}

		// **Correctly build the settings with a new client**
		newSettings := map[string]any{
			"decryption": "none",
			"fallbacks":  []any{},
		}
		
		uuid, err := newUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID for client: %w", err)
		}

		var clientData map[string]any
		if opts.Client != nil {
			// Use provided client options
			clientData = map[string]any{
				"id":         uuid,
				"email":      opts.Client.Email,
				"enable":     opts.Client.Enable,
				"totalGB":    opts.Client.TotalGB,
				"expiryTime": opts.Client.ExpiryTime,
				"limitIp":    opts.Client.LimitIP,
				"flow":       opts.Client.Flow,
				"subId":      opts.Client.SubID,
				"comment":    opts.Client.Comment,
				"reset":      0,
			}
		} else {
			// Create a default client
			clientData = map[string]any{
				"id":      uuid,
				"email":   randSlug(8),
				"enable":  true,
				"totalGB": 0,
				"expiryTime": 0,
				"limitIp": 0,
				"flow":    "",
				"subId":   randSlug(16),
				"reset":   0,
			}
		}
		newSettings["clients"] = []any{clientData}

		v.Set("settings", toJSONString(newSettings))
		v.Set("streamSettings", toJSONString(orig["streamSettings"]))
		v.Set("sniffing", toJSONString(orig["sniffing"]))
		
		// Handle potential nil allocate
		if orig["allocate"] != nil {
			v.Set("allocate", toJSONString(orig["allocate"]))
		} else {
			v.Set("allocate", "{}")
		}

		return v, nil
	}

	const maxTry = 5
	var created xdto.Inbound

	for attempt := 0; attempt < maxTry; attempt++ {
		form, err := buildForm(newRemark, port)
		if err != nil {
			return xdto.Inbound{}, err
		}

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.ep("addInbound", "/panel/api/inbounds/add"),
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var resp map[string]any
		err = d.doJSON(ctx, req, &resp)

		if err != nil {
			if he, ok := err.(*core.HTTPError); ok && he.Code == http.StatusConflict {
				port = 20000 + mrand.Intn(40000) // retry with new port on conflict
				continue
			}
			return xdto.Inbound{}, err // other errors
		}
		
		obj := takeObj(resp)
		if id := anyToInt(obj["id"]); id > 0 {
			created = xdto.Inbound{ID: id, Remark: newRemark, Port: port, Raw: obj}
			break
		}

		// If ID is not in response, try to resolve it by listing inbounds
		for i := 0; i < 5 && created.ID == 0; i++ {
			time.Sleep(time.Duration(200*(i+1)) * time.Millisecond) // exponential backoff
			if got, ok := d.resolveNewInbound(ctx, beforeIDs, newRemark, port); ok {
				created = got
				break
			}
		}

		if created.ID != 0 {
			break
		}
		
		port = 20000 + mrand.Intn(40000)
	}

	if created.ID == 0 {
		return xdto.Inbound{}, fmt.Errorf("clone created but id not resolved (remark=%s, port=%d)", newRemark, port)
	}

	return created, nil
}

// Converts a value to a JSON string.
func toJSONString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case map[string]any, []any:
		b, _ := json.Marshal(t)
		return string(b)
	case nil:
		return "{}" // return empty JSON object for nil
	default:
		// Fallback for other types
		b, _ := json.Marshal(t)
		return string(b)
	}
}