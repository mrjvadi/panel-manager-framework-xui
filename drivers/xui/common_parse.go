package xui

import (
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
)

func getMap(m map[string]any, key string) (map[string]any, bool) {
	if mm, ok := m[key].(map[string]any); ok {
		return mm, true
	}
	return nil, false
}

func inboundFromAny(a any) (xdto.Inbound, bool) {
	switch t := a.(type) {
	case map[string]any:
		if obj, ok := getMap(t, "obj"); ok { // پاسخ‌های { success, msg, obj: {...} }
			return inboundFromAny(obj)
		}
		out := xdto.Inbound{Raw: t}
		if id, ok := t["id"]; ok {
			out.ID = anyToInt(id)
		}
		if p, ok := t["port"]; ok {
			out.Port = anyToInt(p)
		}
		if s, ok := t["remark"].(string); ok {
			out.Remark = s
		}
		if s, ok := t["protocol"].(string); ok {
			out.Protocol = s
		}
		if mm, ok := t["settings"].(map[string]any); ok {
			out.Settings = mm
		}
		if mm, ok := t["sniffing"].(map[string]any); ok {
			out.Sniffing = mm
		}
		if mm, ok := t["streamSettings"].(map[string]any); ok {
			out.StreamSettings = mm
		}
		return out, (out.ID != 0 || out.Port != 0 || out.Remark != "")
	case []any:
		if len(t) == 0 {
			return xdto.Inbound{}, false
		}
		return inboundFromAny(t[0])
	default:
		return xdto.Inbound{}, false
	}
}

func listFromAny(a any) []xdto.Inbound {
	out := []xdto.Inbound{}
	switch t := a.(type) {
	case []any:
		for _, e := range t {
			if inb, ok := inboundFromAny(e); ok {
				out = append(out, inb)
			}
		}
	case map[string]any:
		if inb, ok := inboundFromAny(t); ok {
			out = append(out, inb)
		}
	}
	return out
}
