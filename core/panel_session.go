package core

import "time"
import "context"

type Panel struct{ m *Manager; id string; req *Req }

func (m *Manager) PanelSession(id string, opts ...CtxOption) *Panel { return &Panel{ m:m, id:id, req: m.Request(opts...) } }
func (p *Panel) WithTimeout(d time.Duration) *Panel { p.req.o.timeout = d; return p }
func (p *Panel) Users() ([]User, error) { ctx, cancel := p.req.derive(); defer cancel(); return p.m.Users(ctx, p.id) }
func (p *Panel) Inbounds() ([]Inbound, error) { ctx, cancel := p.req.derive(); defer cancel(); return p.m.Inbounds(ctx, p.id) }

// Generic helpers moved to free functions: As, Try (see as.go)
func (p *Panel) ID() string { return p.id }
func (p *Panel) Manager() *Manager { return p.m }

// TryExt runs fn with a typed extension if available; otherwise returns ErrExtNotSupported.
func TryExt[T any](p *Panel, fn func(ctx context.Context, t T) error) error {
    v, ok := As[T](p.m, p.id); if !ok { return ErrExtNotSupported }
    ctx, cancel := p.req.derive(); defer cancel()
    return fn(ctx, v)
}
