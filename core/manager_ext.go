package core

import (
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

// As: میان‌بُر جنریک برای گرفتن افزونه از یک پنل (بدون نیاز به Ext[T])
func (m *Manager) As[T any](id string) (T, bool) {
	var zero T
	slot, ok := m.getSlot(id)
    if !ok || !slot.enabled {
		return zero, false
	}
	v, ok := any(slot.drv).(T)
    if !ok {
		return zero, false
	}
	return v, true
}

// Panel: API روان‌تر به‌صورت fluent
type PanelHandle struct {
	id string
    m  *Manager
}

func (m *Manager) Panel(id string) PanelHandle { return PanelHandle{id: id, m: m} }

// === هاب افزونه‌ها به‌صورت Fluent ===
// ExtHub: اجازه می‌دهد mgr.ExtHub("id").Marzban() بنویسید (برای جلوگیری از تداخل نام با متد جنریک)

type ExtHub struct{ h PanelHandle }

func (m *Manager) ExtHub(id string) ExtHub { return ExtHub{h: PanelHandle{id: id, m: m}} }

// متدهای تایپ‌شده‌ی عمومی روی PanelHandle
func (h PanelHandle) Usage() (ext.Usage, bool)               { return h.m.As[ext.Usage](h.id) }
func (h PanelHandle) Subscription() (ext.Subscription, bool) { return h.m.As[ext.Subscription](h.id) }

// متدهای مخصوص هر خانواده‌ی درایور روی ExtHub
func (x ExtHub) Marzban() (ext.Marzban, bool) { return x.h.m.As[ext.Marzban](x.h.id) }
func (x ExtHub) XUI() (ext.XUI, bool)         { return x.h.m.As[ext.XUI](x.h.id) }

// نسخه‌های Must* (بدون panic)
func (m *Manager) MustAs[T any](id string) (T, error) {
	if v, ok := m.As[T](id); ok {
		return v, nil
	}
	var zero T
    return zero, ErrExtNotSupported
}

// Try: الگوی اجرایی تمیز برای یک افزونه
func (m *Manager) Try[T any](id string, fn func(T) error) error {
	v, ok := m.As[T](id)
    if !ok {
		return ErrExtNotSupported
	}
	return fn(v)
}
