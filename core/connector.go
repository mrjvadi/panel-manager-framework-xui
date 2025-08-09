package core

import (
    "context"
    "sync"
)

// Connector: واسط اختیاری برای آماده‌سازی اتصال (مثلاً Login اولیه)
type Connector interface {
    Connect(ctx context.Context) error
}

func (m *Manager) Connect(id string, ctx context.Context) error {
    if s, ok := m.getSlot(id); ok && s.enabled {
        if c, ok := s.drv.(Connector); ok {
            return c.Connect(ctx)
        }
    }
    return nil
}

func (m *Manager) ConnectAll(ctx context.Context) error {
    ids, slots := m.snapshotEnabled()
    sem := make(chan struct{}, m.opts.MaxConcurrency)
    errc := make(chan error, len(ids))
    var wg sync.WaitGroup
    for i, id := range ids {
        wg.Add(1); sem <- struct{}{}; slot := slots[i]; panelID := id
        go func(id string, slot *panelSlot) {
            defer wg.Done(); defer func(){<-sem}()
            if c, ok := slot.drv.(Connector); ok {
                if err := c.Connect(ctx); err != nil { errc <- err }
            }
        }(id, slot)
    }
    wg.Wait(); close(errc)
    for e := range errc { if e != nil { return e } }
    return nil
}
