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

// اطمینان از اینکه درایور، این اکستنشن را هم اکسپوز می‌کند:
var _ ext.XUIShallowClone = (*sanaei)(nil)

// ---- helpers (لوکال و امن) ----

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

// از ساختارهای متنوع پاسخ‌ها، آبجکت اصلی را بیرون می‌کشد
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

func randSlug(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	for i := 0; i < n; i++ {
		ri, _ := crand.Int(crand.Reader, big.NewInt(int64(len(letters))))
		b.WriteByte(letters[ri.Int64()])
	}
	return b.String()
}

func (d *sanaei) addClientFromCreate(ctx context.Context, inboundID int, c xdto.ClientCreate) error {
	if c.Email == "" {
		return fmt.Errorf("client email is required")
	}
	cli := map[string]any{
		"email":      c.Email,
		"enable":     c.Enable,
		"expiryTime": c.ExpiryTime,
		"totalGB":    c.TotalGB,
		"limitIp":    c.LimitIP,
		"reset":      0,
		"flow":       c.Flow,
		"subId":      c.SubID,
		"comment":    c.Comment,
	}
	settings := map[string]any{"clients": []any{cli}}
	sb, _ := json.Marshal(settings)

	payload := map[string]any{
		"id":       inboundID,
		"settings": string(sb), // اکثر فورک‌ها string می‌خوان
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.ep("addClient", "/panel/api/inbounds/addClient"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return d.doJSON(ctx, req, nil)
}

// ---- پیاده‌سازی اصلی: CloneInboundShallow ----
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

	// قبل از add: شناسه‌های موجود برای resolve
	before, _ := d.listInboundsRaw(ctx)
	beforeIDs := make(map[int]struct{}, len(before))
	for _, m := range before {
		if id := anyToInt(m["id"]); id > 0 {
			beforeIDs[id] = struct{}{}
		}
	}

	// بدنهٔ form-urlencoded دقیقاً مثل نمونهٔ پنل
	buildForm := func(rem string, p int) url.Values {
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
		v.Set("settings", toJSONString(orig["settings"]))
		v.Set("streamSettings", toJSONString(orig["streamSettings"]))
		v.Set("sniffing", toJSONString(orig["sniffing"]))
		v.Set("allocate", toJSONString(orig["allocate"]))
		return v
	}

	const maxTry = 5
	var created xdto.Inbound

	for attempt := 0; attempt < maxTry; attempt++ {
		form := buildForm(newRemark, port)
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+d.ep("addInbound", "/panel/api/inbounds/add"),
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var resp map[string]any
		err := d.doJSON(ctx, req, &resp)
		if err != nil {
			// اگر پورت تکراری بود، یه پورت جدید انتخاب کن و دوباره
			if he, ok := err.(*core.HTTPError); ok && he.Code == http.StatusConflict {
				port = 20000 + mrand.Intn(40000)
				continue
			}
			return xdto.Inbound{}, err
		}

		// بعضی نسخه‌ها id را در obj می‌دهند، بعضی‌ها نه
		obj := takeObj(resp)
		if id := anyToInt(obj["id"]); id > 0 {
			created = xdto.Inbound{ID: id, Remark: newRemark, Port: port, Raw: obj}
			break
		}

		// اگر id نگرفتیم، کمی صبر و resolve از list
		for i := 0; i < 5 && created.ID == 0; i++ {
			time.Sleep(150 * time.Millisecond)
			if got, ok := d.resolveNewInbound(ctx, beforeIDs, newRemark, port); ok {
				created = got
				break
			}
		}

		if created.ID != 0 {
			break
		}

		// اگر هنوز پیدا نشد، پورت رو عوض کنیم شاید conflict بی‌سروصدا رخ داده بوده
		port = 20000 + mrand.Intn(40000)
	}

	if created.ID == 0 {
		return xdto.Inbound{}, fmt.Errorf("clone created but id not resolved (remark=%s, port=%d)", newRemark, port)
	}

	// اگر Caller کلاینت داده، مطابق همان مقادیر خودت بسازیم
	if opts.Client != nil {
		if err := d.addClientFromCreate(ctx, created.ID, *opts.Client); err != nil {
			return created, fmt.Errorf("inbound cloned, add-client failed: %w", err)
		}
	}

	return created, nil
}

// داخل drivers/xui/sanaei_shallow.go، نزدیک بقیه هلسپرها
func toJSONString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case map[string]any, []any:
		b, _ := json.Marshal(t)
		return string(b)
	case nil:
		return ""
	default:
		// اگر نوع دیگری آمد، آخرین تلاش: marshal
		b, _ := json.Marshal(t)
		return string(b)
	}
}
