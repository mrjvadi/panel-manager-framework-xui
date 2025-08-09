package xui

import (
	"context"
	"fmt"
	"github.com/mrjvadi/panel-manager-framework-xui/core"
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
	"math/rand"
	"net/http"
	"time"
)

const SanaeiName = "xui.sanaei"
const PluginVer = "v2.6.1"

func init() { core.Register(SanaeiName, NewSanaei) }

type sanaei struct{ *generic }

func NewSanaei(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
	g := newGeneric(sp, opts...)
	return &sanaei{g}, nil
}

func (d *sanaei) Name() string { return SanaeiName }
func (d *sanaei) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	return PluginVer
}
func (d *sanaei) Capabilities() core.Capabilities { return core.Capabilities{} }

// XUIDeleter
func (d *sanaei) DeleteInboundByID(ctx context.Context, id int) error {
	path := fmt.Sprintf(d.ep("deleteInbound", "/panel/inbound/del/%d"), id)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, d.cli.BaseURL+path, nil)
	return d.doJSON(ctx, req, nil)
}

func (d *sanaei) ListUsers(ctx context.Context) ([]core.User, error) { return nil, nil }

func (d *sanaei) ListInbounds(ctx context.Context) ([]core.Inbound, error) {
	var body any
	if err := d.getJSON(ctx, d.sp.Endpoints["listInbounds"], &body); err != nil {
		return nil, err
	}
	return nil, nil
}

// Typed API
var _ ext.XUITyped = (*sanaei)(nil)

func (d *sanaei) AddInboundTyped(ctx context.Context, in xdto.InboundCreate) (xdto.Inbound, error) {
	body := map[string]any{
		"remark":         in.Remark,
		"protocol":       in.Protocol,
		"port":           in.Port,
		"settings":       in.Settings,
		"sniffing":       in.Sniffing,
		"streamSettings": in.StreamSettings,
	}
	var resp map[string]any
	if err := d.postJSON(ctx, d.sp.Endpoints["addInbound"], body, &resp); err != nil {
		return xdto.Inbound{}, err
	}
	// اگر خود پاسخ id داد، همونو برگردون
	if out, ok := inboundFromAny(resp); ok && out.ID != 0 {
		return out, nil
	}
	// وگرنه بده به کالر تا خودش resolve کنه
	return xdto.Inbound{Raw: resp}, nil
}

func (d *sanaei) UpdateInboundTyped(ctx context.Context, inboundID int, in xdto.InboundUpdate) (xdto.Inbound, error) {
	path := fmt.Sprintf(d.sp.Endpoints["updateInbound"], inboundID)
	var resp map[string]any
	body := map[string]any{
		"remark":         in.Remark,
		"protocol":       in.Protocol,
		"port":           in.Port,
		"settings":       in.Settings,
		"sniffing":       in.Sniffing,
		"streamSettings": in.StreamSettings,
	}
	if err := d.postJSON(ctx, path, body, &resp); err != nil {
		return xdto.Inbound{}, err
	}
	if out, ok := inboundFromAny(resp); ok && out.ID != 0 {
		return out, nil
	}
	return d.GetInboundTyped(ctx, inboundID)
}

func (d *sanaei) ClientTrafficByEmailTyped(ctx context.Context, email string) ([]xdto.TrafficRecord, error) {
	return nil, nil
}
func (d *sanaei) CloneInboundTyped(ctx context.Context, inboundID int, opts xdto.CloneInboundOptions) (xdto.Inbound, error) {
	orig, err := d.GetInboundTyped(ctx, inboundID)
	if err != nil {
		return xdto.Inbound{}, err
	}

	// baseline
	baseRemark := orig.Remark
	if baseRemark == "" {
		baseRemark = "inb"
	}

	// remark هدف
	newRemark := baseRemark + "-copy-" + time.Now().Format("20060102-150405")
	if opts.Remark != nil && *opts.Remark != "" {
		newRemark = *opts.Remark
	}

	// پورت هدف
	port := 0
	if opts.Port != nil {
		port = *opts.Port
	} else {
		port = 20000 + rand.Intn(40000)
	}

	// قبل از add: لیست شناسه‌های موجود
	before, _ := d.ListInboundsTyped(ctx)
	beforeSet := make(map[int]struct{}, len(before))
	for _, it := range before {
		if it.ID > 0 {
			beforeSet[it.ID] = struct{}{}
		}
	}

	const maxTry = 5
	var created xdto.Inbound
	for i := 0; i < maxTry; i++ {
		in := xdto.InboundCreate{
			Remark:         newRemark,
			Protocol:       orig.Protocol,
			Port:           port,
			Settings:       orig.Settings,
			Sniffing:       orig.Sniffing,
			StreamSettings: orig.StreamSettings,
		}
		out, err := d.AddInboundTyped(ctx, in)
		if err == nil && out.ID != 0 {
			created = out
			break
		}
		if err == nil && out.ID == 0 {
			// پاسخ id نداد؛ resolve کن
			if got, ok := d.resolveNewInbound(ctx, beforeSet, newRemark, port); ok {
				created = got
				break
			}
			// ادامه بده: شاید conflict بعدی باشه
		}
		// اگر کانفلیکت پورت بود، پورت جدید امتحان کن
		if he, ok := err.(*core.HTTPError); ok && he.Code == http.StatusConflict {
			port = 20000 + rand.Intn(40000)
			continue
		}
		// سایر خطاها
		if err != nil {
			return xdto.Inbound{}, err
		}
	}

	if created.ID == 0 {
		// آخرین تلاش برای resolve
		if got, ok := d.resolveNewInbound(ctx, beforeSet, newRemark, port); ok {
			created = got
		}
	}

	if created.ID == 0 {
		return xdto.Inbound{}, fmt.Errorf("clone created but id not resolved (remark=%s, port=%d)", newRemark, port)
	}
	return created, nil
}

func (d *sanaei) ListInboundsTyped(ctx context.Context) ([]xdto.Inbound, error) {
	var resp any
	if err := d.getJSON(ctx, d.sp.Endpoints["listInbounds"], &resp); err != nil {
		return nil, err
	}
	switch v := resp.(type) {
	case map[string]any:
		if obj, ok := v["obj"]; ok {
			return listFromAny(obj), nil
		}
		return listFromAny(v), nil
	default:
		return listFromAny(v), nil
	}
}

func (d *sanaei) GetInboundTyped(ctx context.Context, inboundID int) (xdto.Inbound, error) {
	path := fmt.Sprintf(d.sp.Endpoints["getInbound"], inboundID)
	fmt.Println(path)
	var resp map[string]any
	if err := d.getJSON(ctx, path, &resp); err != nil {
		return xdto.Inbound{}, err
	}
	if out, ok := inboundFromAny(resp); ok {
		return out, nil
	}
	return xdto.Inbound{Raw: resp}, nil
}
