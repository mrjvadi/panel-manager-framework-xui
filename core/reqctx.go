package core

import (
    "context"
    "time"
)

// CtxOption: تنظیمات محدوده‌ی درخواست
type CtxOption func(*reqopts)

type reqopts struct {
    timeout time.Duration
    vals    []ctxKV
}

type ctxKV struct{ k, v any }

// WithReqTimeoutOpt: تایم‌اوت پیش‌فرض برای هر فراخوانی در این اسکوپ
func WithReqTimeoutOpt(d time.Duration) CtxOption { return func(o *reqopts) { if d > 0 { o.timeout = d } } }

// WithValueOpt: تزریق مقدار به context
func WithValueOpt(k, v any) CtxOption { return func(o *reqopts) { o.vals = append(o.vals, ctxKV{k:k, v:v}) } }

// Req: یک اسکوپ کانتکست‌دار برای اجرای فراخوانی‌ها بدون ساخت context در هر بار
type Req struct {
    m   *Manager
    base context.Context
    o   reqopts
}

// Request: ساخت اسکوپ بر پایه‌ی BaseCtx مدیر و ReqTimeout پیش‌فرض
func (m *Manager) Request(opts ...CtxOption) *Req {
    r := &Req{ m: m, base: m.opts.BaseCtx, o: reqopts{ timeout: m.opts.ReqTimeout } }
    for _, fn := range opts { fn(&r.o) }
    return r
}

// RequestFrom: ساخت اسکوپ با کانتکست پایه‌ی سفارشی (مثلاً Context والد سرویس شما)
func (m *Manager) RequestFrom(base context.Context, opts ...CtxOption) *Req {
    if base == nil { base = m.opts.BaseCtx }
    r := &Req{ m: m, base: base, o: reqopts{ timeout: m.opts.ReqTimeout } }
    for _, fn := range opts { fn(&r.o) }
    return r
}

// derive: ساخت کانتکست فرزند با رعایت timeout و Valueها
func (r *Req) derive() (context.Context, context.CancelFunc) {
    ctx := r.base
    var cancel context.CancelFunc = func() {}
    if r.o.timeout > 0 {
        ctx, cancel = context.WithTimeout(ctx, r.o.timeout)
    } else {
        ctx, cancel = context.WithCancel(ctx)
    }
    for _, kv := range r.o.vals {
        ctx = context.WithValue(ctx, kv.k, kv.v)
    }
    return ctx, cancel
}

// WithTimeout: تغییر تایم‌اوت برای این اسکوپ (chainable)
func (r *Req) WithTimeout(d time.Duration) *Req { r.o.timeout = d; return r }
// WithValue: افزودن Value به کانتکست مشتقات (chainable)
func (r *Req) WithValue(k, v any) *Req { r.o.vals = append(r.o.vals, ctxKV{k:k, v:v}); return r }

// Do: اجرای یک تابع با کانتکست مشتق‌شده از اسکوپ
func (r *Req) Do(fn func(ctx context.Context) error) error {
    ctx, cancel := r.derive(); defer cancel()
    return fn(ctx)
}

// === شورتکات‌های متداول ===

// UsersAll: اجرای Manager.UsersAll با کانتکست داخلی
func (r *Req) UsersAll() (map[string][]User, error) {
    ctx, cancel := r.derive(); defer cancel()
    return r.m.UsersAll(ctx)
}

func (r *Req) InboundsAll() (map[string][]Inbound, error) {
    ctx, cancel := r.derive(); defer cancel()
    return r.m.InboundsAll(ctx)
}

func (r *Req) Users(id string) ([]User, error) {
    ctx, cancel := r.derive(); defer cancel()
    return r.m.Users(ctx, id)
}

func (r *Req) Inbounds(id string) ([]Inbound, error) {
    ctx, cancel := r.derive(); defer cancel()
    return r.m.Inbounds(ctx, id)
}

// GroupCtx: نسخه‌ی کانتکست‌دار از Group برای اجرای عملیات گروهی بدون ساخت ctx
type GroupCtx struct {
    r *Req
    g Group
}

func (r *Req) Group(g Group) GroupCtx { return GroupCtx{ r: r, g: g } }

// میانبرهای خانواده‌ها
func (r *Req) Marzban() GroupCtx    { return r.Group(r.m.Marzban()) }
func (r *Req) XUIAll() GroupCtx     { return r.Group(r.m.XUIAll()) }
func (r *Req) XUIGeneric() GroupCtx { return r.Group(r.m.XUIGeneric()) }
func (r *Req) XUISanaei() GroupCtx  { return r.Group(r.m.XUISanaei()) }
func (r *Req) XUIAlireza() GroupCtx { return r.Group(r.m.XUIAlireza()) }

func (g GroupCtx) IDs() []string { return g.g.IDs() }

func (g GroupCtx) UsersAll() (map[string][]User, error) {
    ctx, cancel := g.r.derive(); defer cancel()
    return g.g.UsersAll(ctx)
}

func (g GroupCtx) InboundsAll() (map[string][]Inbound, error) {
    ctx, cancel := g.r.derive(); defer cancel()
    return g.g.InboundsAll(ctx)
}

// TryEachCtx: اجرای تابع روی همه‌ی پنل‌های گروه با تزریق ctx
func (g GroupCtx) TryEachCtx[T any](fn func(ctx context.Context, id string, t T) error) error {
    items := g.g.ExtMap[T]()
    // اجرای موازی با همزمانی Manager
    sem := make(chan struct{}, g.r.m.opts.MaxConcurrency)
    errc := make(chan error, len(items))
    done := make(chan struct{})
    go func() {
        for id, v := range items {
            sem <- struct{}{} // acquire
            id, v := id, v
            go func() {
                defer func(){ <-sem }()
                ctx, cancel := g.r.derive(); defer cancel()
                if err := fn(ctx, id, v); err != nil { errc <- err }
            }()
        }
        // drain semaphore
        for i := 0; i < cap(sem); i++ { sem <- struct{}{} }
        close(done)
    }()
    <-done
    close(errc)
    for e := range errc { if e != nil { return e } }
    return nil
}
