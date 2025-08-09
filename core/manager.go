package core

import (
    "context"
    "errors"
    "log"
    "log/slog"
    "net/http"
    "sync"
    "time"
)

type Option func(*options)

type options struct {
    HTTPClient      *http.Client
    Timeout         time.Duration
    Hooks           *Hooks
    MaxConcurrency  int
    RetryPolicy     RetryPolicy
    BreakerThresh   int
    BreakerCooldown time.Duration
    BaseCtx         context.Context
    ReqTimeout      time.Duration
    Logger          Logger
}

func WithHTTPClient(c *http.Client) Option     { return func(o *options) { o.HTTPClient = c } }
func WithTimeout(d time.Duration) Option       { return func(o *options) { if d > 0 { o.Timeout = d } } }
func WithHooks(h *Hooks) Option                { return func(o *options) { o.Hooks = h } }
func WithMaxConcurrency(n int) Option          { return func(o *options) { if n > 0 { o.MaxConcurrency = n } } }
func WithRetryPolicy(p RetryPolicy) Option     { return func(o *options) { o.RetryPolicy = p } }
func WithBreaker(th int, cd time.Duration) Option {
    return func(o *options) {
        if th > 0 { o.BreakerThresh = th }
        if cd > 0 { o.BreakerCooldown = cd }
    }
}
func WithBaseContext(ctx context.Context) Option { return func(o *options) { if ctx != nil { o.BaseCtx = ctx } } }
func WithRequestTimeout(d time.Duration) Option  { return func(o *options) { if d > 0 { o.ReqTimeout = d } } }
func WithLogger(l Logger) Option                 { return func(o *options) { if l != nil { o.Logger = l } } }
func WithSlogLogger(l *slog.Logger) Option       { return func(o *options) { if l != nil { o.Logger = NewSlogAdapter(l) } } }
func WithStdLogger(l *log.Logger) Option         { return func(o *options) { if l != nil { o.Logger = NewStdLoggerAdapter(l) } } }

type Manager struct {
    mu sync.RWMutex
    panels map[string]*panelSlot
    opts options
}

type panelSlot struct {
    drv Driver
    enabled bool
}

func New(opts ...Option) *Manager {
    o := options{ Timeout: 30*time.Second, MaxConcurrency: 8, BaseCtx: context.Background(), ReqTimeout: 15*time.Second, Logger: NoopLogger() }
    for _, fn := range opts { fn(&o) }
    return &Manager{ panels: map[string]*panelSlot{}, opts: o }
}

func (m *Manager) Attach(spec PanelSpec, driverName string, opts ...Option) error {
    inherit := []Option{ WithHTTPClient(m.opts.HTTPClient), WithTimeout(m.opts.Timeout), WithRetryPolicy(m.opts.RetryPolicy), WithBreaker(m.opts.BreakerThresh, m.opts.BreakerCooldown), WithLogger(m.opts.Logger) }
    opts = append(inherit, opts...)
    f, ok := Factory(driverName); if !ok { return errors.New("driver not found: "+driverName) }
    drv, err := f(spec, opts...); if err != nil { return err }
    m.mu.Lock(); m.panels[spec.ID] = &panelSlot{ drv: drv, enabled: true }; m.mu.Unlock()
    return nil
}
func (m *Manager) AttachByKind(spec PanelSpec, kind DriverKind, opts ...Option) error { return m.Attach(spec, kind.String(), opts...) }
func (m *Manager) Detach(id string) { m.mu.Lock(); delete(m.panels, id); m.mu.Unlock() }
func (m *Manager) Enable(id string) { m.mu.Lock(); if s, ok := m.panels[id]; ok { s.enabled = true }; m.mu.Unlock() }
func (m *Manager) Disable(id string){ m.mu.Lock(); if s, ok := m.panels[id]; ok { s.enabled = false }; m.mu.Unlock() }
func (m *Manager) getSlot(id string) (*panelSlot, bool) { m.mu.RLock(); defer m.mu.RUnlock(); s, ok := m.panels[id]; return s, ok }

func (m *Manager) snapshotEnabled() ([]string, []*panelSlot) {
    m.mu.RLock(); defer m.mu.RUnlock()
    ids := make([]string,0,len(m.panels)); slots := make([]*panelSlot,0,len(m.panels))
    for id, s := range m.panels { if s.enabled { ids = append(ids, id); slots = append(slots, s) } }
    return ids, slots
}

func (m *Manager) UsersAll(ctx context.Context) (map[string][]User, error) {
    out := map[string][]User{}; ids, slots := m.snapshotEnabled()
    sem := make(chan struct{}, m.opts.MaxConcurrency); var wg sync.WaitGroup; var mu sync.Mutex
    for i, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string, s *panelSlot){ defer wg.Done(); defer func(){<-sem}(); rows, err := s.drv.ListUsers(ctx); if err != nil { rows = nil }; mu.Lock(); out[id]=rows; mu.Unlock() }(id, slots[i])
    }
    wg.Wait(); return out, nil
}
func (m *Manager) InboundsAll(ctx context.Context) (map[string][]Inbound, error) {
    out := map[string][]Inbound{}; ids, slots := m.snapshotEnabled()
    sem := make(chan struct{}, m.opts.MaxConcurrency); var wg sync.WaitGroup; var mu sync.Mutex
    for i, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string, s *panelSlot){ defer wg.Done(); defer func(){<-sem}(); rows, err := s.drv.ListInbounds(ctx); if err != nil { rows = nil }; mu.Lock(); out[id]=rows; mu.Unlock() }(id, slots[i])
    }
    wg.Wait(); return out, nil
}
func (m *Manager) Users(ctx context.Context, id string) ([]User, error) { s, ok := m.getSlot(id); if !ok || !s.enabled { return nil, errors.New("panel not found or disabled") }; return s.drv.ListUsers(ctx) }
func (m *Manager) Inbounds(ctx context.Context, id string) ([]Inbound, error) { s, ok := m.getSlot(id); if !ok || !s.enabled { return nil, errors.New("panel not found or disabled") }; return s.drv.ListInbounds(ctx) }

func (m *Manager) As[T any](id string) (T, bool) { s, ok := m.getSlot(id); if !ok || !s.enabled { var zero T; return zero, false }; v, ok := any(s.drv).(T); return v, ok }
