package core

import "context"
import "time"

type CtxOption func(*reqopts)
type reqopts struct { timeout time.Duration }
func WithReqTimeoutOpt(d time.Duration) CtxOption { return func(o *reqopts){ if d>0 { o.timeout=d } } }

type Req struct { m *Manager; base context.Context; o reqopts }
func (m *Manager) Request(opts ...CtxOption) *Req {
    r := &Req{ m:m, base: m.opts.BaseCtx, o: reqopts{ timeout: m.opts.ReqTimeout } }
    for _, fn := range opts { fn(&r.o) }
    return r
}
func (r *Req) derive() (context.Context, context.CancelFunc) {
    if r.o.timeout > 0 { return context.WithTimeout(r.base, r.o.timeout) }
    return context.WithCancel(r.base)
}
func (r *Req) UsersAll() (map[string][]User, error) { ctx, cancel := r.derive(); defer cancel(); return r.m.UsersAll(ctx) }
