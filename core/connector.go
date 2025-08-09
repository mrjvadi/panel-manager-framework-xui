package core

import (
	"context"
	"fmt"
	"time"
)

// Connector یک هندل سبک‌وزن برای کار با یک پنل خاص است.
// مزیتش اینه که: ۱) همیشه با همون panelID کار می‌کنه
// ۲) خودش کانتکست داخلی با timeout از Manager می‌سازه
// ۳) می‌تونی به‌صورت type-safe به درایورهای توسعه‌یافته دسترسی بگیری (با As[T] آزاد)
type Connector struct {
	m       *Manager
	panelID string
	req     *Req
}

// ساخت Connector روی یک panelID مشخص
func (m *Manager) Connector(panelID string, opts ...CtxOption) *Connector {
	return &Connector{
		m:       m,
		panelID: panelID,
		req:     m.Request(opts...),
	}
}

// تنظیم timeout اختصاصی برای همهٔ درخواست‌های این Connector
func (c *Connector) WithTimeout(d time.Duration) *Connector {
	c.req.o.timeout = d
	return c
}

// گرفتن Panel session هم‌مسیر با همین connector (برای استفاده از شورتکات‌های Panel)
func (c *Connector) Panel() *Panel {
	return &Panel{m: c.m, id: c.panelID, req: c.req}
}

// اطمینان از دسترسی به پنل و انجام Login (در صورت نیاز)
func (c *Connector) EnsureLogin() error {
	s, ok := c.m.getSlot(c.panelID)
	if !ok || !s.enabled {
		return fmt.Errorf("panel not found or disabled: %s", c.panelID)
	}
	ctx, cancel := c.req.derive()
	defer cancel()
	return s.drv.Login(ctx)
}

// لیست کاربران پنل هدف
func (c *Connector) Users() ([]User, error) {
	ctx, cancel := c.req.derive()
	defer cancel()
	return c.m.Users(ctx, c.panelID)
}

// لیست این‌باندهای پنل هدف
func (c *Connector) Inbounds() ([]Inbound, error) {
	ctx, cancel := c.req.derive()
	defer cancel()
	return c.m.Inbounds(ctx, c.panelID)
}

// اجرای تابع روی یک اکستنشن/درایور تایپ‌شده اگر موجود باشد.
// دقت کن: این تابع از free-function جنریک As[T] استفاده می‌کند، نه متد جنریک روی رسیور.
// تابع آزاد جنریک: به‌جای متد جنریک روی Connector
func UseExt[T any](c *Connector, fn func(ctx context.Context, drv T) error) error {
	drv, ok := As[T](c.m, c.panelID)
	if !ok {
		return ErrExtNotSupported
	}
	ctx, cancel := c.req.derive()
	defer cancel()
	return fn(ctx, drv)
}
