package core

import (
    "context"
    "time"
)

// Panel: یک wrapper روی یک پنل خاص با کانتکست داخلی
type Panel struct {
    m   *Manager
    id  string
    req *Req
}

// PanelSession: ساخت wrapper پنل با کانتکست داخلی
func (m *Manager) PanelSession(id string, opts ...CtxOption) *Panel {
    return &Panel{ m: m, id: id, req: m.Request(opts...) }
}

// Chainable تنظیمات
func (p *Panel) WithTimeout(d time.Duration) *Panel { p.req.WithTimeout(d); return p }
func (p *Panel) WithValue(k, v any) *Panel          { p.req.WithValue(k, v); return p }

// دسترسی به درایور خام در صورت نیاز
func (p *Panel) driver() (Driver, bool) {
    s, ok := p.m.getSlot(p.id)
    if !ok || !s.enabled { return nil, false }
    return s.drv, true
}

// شورتکات‌های رایج بدون ctx بیرونی
func (p *Panel) Users() ([]User, error) {
    ctx, cancel := p.req.derive(); defer cancel()
    return p.m.Users(ctx, p.id)
}
func (p *Panel) Inbounds() ([]Inbound, error) {
    ctx, cancel := p.req.derive(); defer cancel()
    return p.m.Inbounds(ctx, p.id)
}

// As: گرفتن اکستنشن برای همین پنل
func (p *Panel) As[T any]() (T, bool) {
    return p.m.As[T](p.id)
}

// Try: اجرای تابع روی یک اکستنشن برای همین پنل با ctx داخلی
func (p *Panel) Try[T any](fn func(ctx context.Context, t T) error) error {
    v, ok := p.m.As[T](p.id); if !ok { return ErrExtNotSupported }
    ctx, cancel := p.req.derive(); defer cancel()
    return fn(ctx, v)
}

// ==== خانواده‌های رایج به صورت wrapper ====

type PanelMarzban struct {
    p  *Panel
    mz any // می‌تونه ext.Marzban یا ext.MarzbanTyped باشه
}

func (p *Panel) Marzban() (PanelMarzban, bool) {
    if v, ok := p.m.As[any](p.id); ok {
        // وجود هرکدام کافی‌ست؛ خود متدها تشخیص می‌دهند
        return PanelMarzban{ p: p, mz: v }, true
    }
    return PanelMarzban{}, false
}

type PanelXUI struct { p *Panel; x any }

func (p *Panel) XUI() (PanelXUI, bool) {
    if v, ok := p.m.As[any](p.id); ok {
        return PanelXUI{ p: p, x: v }, true
    }
    return PanelXUI{}, false
}
