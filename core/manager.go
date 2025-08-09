package core

import (
    "context"
    "errors"
    "net/http"
    "sync"
    "time"
)

type Option func(*options)

type options struct {
    HTTPClient     *http.Client
    Timeout        time.Duration
    Hooks          *Hooks
    MaxConcurrency int // حداکثر همزمانی
}

func WithHTTPClient(c *http.Client) Option { return func(o *options) { o.HTTPClient = c } }
func WithTimeout(d time.Duration) Option    { return func(o *options) { o.Timeout = d } }
func WithHooks(h *Hooks) Option             { return func(o *options) { o.Hooks = h } }
func WithMaxConcurrency(n int) Option       { return func(o *options) { if n > 0 { o.MaxConcurrency = n } } }

type panelSlot struct {
    drv     Driver
    enabled bool
}

type Manager struct {
    mu     sync.RWMutex
    opts   options
    panels map[string]*panelSlot // key: PanelSpec.ID
}

func New(opts ...Option) *Manager {
    o := options{ Timeout: 30 * time.Second, MaxConcurrency: 8 }
    for _, fn := range opts { fn(&o) }
    return &Manager{ opts: o, panels: map[string]*panelSlot{} }
}

func (m *Manager) AttachByKind(spec PanelSpec, kind DriverKind, opts ...Option) error { return m.Attach(spec, kind.String(), opts...) }

func (m *Manager) Attach(spec PanelSpec, driverName string, opts ...Option) error {
    m.mu.Lock(); defer m.mu.Unlock()
    f, ok := Factory(driverName)
    if !ok { return errors.New("unknown driver: " + driverName) }
    drv, err := f(spec, opts...)
    if err != nil { return err }
    m.panels[spec.ID] = &panelSlot{ drv: drv, enabled: true }
    return nil
}

func (m *Manager) Detach(id string)                 { m.mu.Lock(); delete(m.panels, id); m.mu.Unlock() }
func (m *Manager) Disable(id string)                { m.mu.Lock(); if s, ok := m.panels[id]; ok { s.enabled = false }; m.mu.Unlock() }
func (m *Manager) Enable(id string)                 { m.mu.Lock(); if s, ok := m.panels[id]; ok { s.enabled = true }; m.mu.Unlock() }
func (m *Manager) getSlot(id string) (*panelSlot, bool) { m.mu.RLock(); defer m.mu.RUnlock(); s, ok := m.panels[id]; return s, ok }

func (m *Manager) snapshotEnabled() ([]string, []*panelSlot) {
    m.mu.RLock(); defer m.mu.RUnlock()
    ids := make([]string, 0, len(m.panels)); slots := make([]*panelSlot, 0, len(m.panels))
    for id, s := range m.panels { if s.enabled { ids = append(ids, id); slots = append(slots, s) } }
    return ids, slots
}

func (m *Manager) UsersAll(ctx context.Context) (map[string][]User, error) {
    ids, slots := m.snapshotEnabled()
    out := make(map[string][]User, len(ids))
    var mu sync.Mutex
    sem := make(chan struct{}, m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for i, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string, slot *panelSlot) { defer wg.Done(); defer func(){<-sem}()
            rows, err := slot.drv.ListUsers(ctx); if err != nil { rows = []User{} }
            mu.Lock(); out[id] = rows; mu.Unlock()
        }(id, slots[i])
    }
    wg.Wait(); return out, nil
}

func (m *Manager) InboundsAll(ctx context.Context) (map[string][]Inbound, error) {
    ids, slots := m.snapshotEnabled()
    out := make(map[string][]Inbound, len(ids))
    var mu sync.Mutex
    sem := make(chan struct{}, m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for i, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string, slot *panelSlot) { defer wg.Done(); defer func(){<-sem}()
            rows, err := slot.drv.ListInbounds(ctx); if err != nil { rows = []Inbound{} }
            mu.Lock(); out[id] = rows; mu.Unlock()
        }(id, slots[i])
    }
    wg.Wait(); return out, nil
}

func (m *Manager) Users(ctx context.Context, id string) ([]User, error)    { slot, ok := m.getSlot(id); if !ok || !slot.enabled { return nil, errors.New("panel not found or disabled") }; return slot.drv.ListUsers(ctx) }
func (m *Manager) Inbounds(ctx context.Context, id string) ([]Inbound, error) { slot, ok := m.getSlot(id); if !ok || !slot.enabled { return nil, errors.New("panel not found or disabled") }; return slot.drv.ListInbounds(ctx) }
