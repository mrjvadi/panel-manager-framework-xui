package core

import (
    "context"
    "strings"
    "sync"
)

type Group struct {
    m      *Manager
    filter func(id string, s *panelSlot) bool
}

// Kind: تمام پنل‌هایی که Name() آن‌ها برابر با DriverKind.String() باشد
func (m *Manager) Kind(kind DriverKind) Group {
    name := kind.String()
    return Group{ m: m, filter: func(_ string, s *panelSlot) bool { return s.enabled && s.drv.Name() == name } }
}

// Family: همه‌ی پنل‌هایی که driver.Name با پیشوند داده‌شده شروع شود (مثلاً "xui.")
func (m *Manager) Family(prefix string) Group {
    return Group{ m: m, filter: func(_ string, s *panelSlot) bool { return s.enabled && strings.HasPrefix(s.drv.Name(), prefix) } }
}

// Where: فیلتر سفارشی
func (m *Manager) Where(pred func(id string, d Driver) bool) Group {
    return Group{ m: m, filter: func(id string, s *panelSlot) bool { return s.enabled && pred(id, s.drv) } }
}

// Filter: پالایش بیشتر روی یک Group موجود (AND)
func (g Group) Filter(pred func(id string, d Driver) bool) Group {
    return Group{ m: g.m, filter: func(id string, s *panelSlot) bool { return g.filter(id, s) && pred(id, s.drv) } }
}

// === فیلترهای ورژنی ===
func (m *Manager) VersionEq(v string) Group        { return m.Where(func(_ string, d Driver) bool { return compareVersionStr(d.Version(), v) == 0 }) }
func (m *Manager) VersionPrefix(p string) Group    { return m.Where(func(_ string, d Driver) bool { return strings.HasPrefix(d.Version(), p) }) }
func (m *Manager) VersionRange(min, max string) Group { return m.Where(func(_ string, d Driver) bool { return compareVersionStr(d.Version(), min) >= 0 && compareVersionStr(d.Version(), max) <= 0 }) }

func (g Group) WhereVersionEq(v string) Group        { return g.Filter(func(_ string, d Driver) bool { return compareVersionStr(d.Version(), v) == 0 }) }
func (g Group) WhereVersionPrefix(p string) Group    { return g.Filter(func(_ string, d Driver) bool { return strings.HasPrefix(d.Version(), p) }) }
func (g Group) WhereVersionRange(min, max string) Group { return g.Filter(func(_ string, d Driver) bool { return compareVersionStr(d.Version(), min) >= 0 && compareVersionStr(d.Version(), max) <= 0 }) }

// IDs: خروجی گرفتن شناسه‌های پنلِ انتخاب‌شده
func (g Group) IDs() []string {
    g.m.mu.RLock(); defer g.m.mu.RUnlock()
    ids := make([]string, 0, len(g.m.panels))
    for id, s := range g.m.panels { if g.filter(id, s) { ids = append(ids, id) } }
    return ids
}

// UsersAll: معادل Manager.UsersAll اما فقط برای گروه فعلی
func (g Group) UsersAll(ctx context.Context) (map[string][]User, error) {
    ids := g.IDs()
    out := make(map[string][]User, len(ids))
    var mu sync.Mutex
    sem := make(chan struct{}, g.m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for _, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string) { defer wg.Done(); defer func(){<-sem}()
            if slot, ok := g.m.getSlot(id); ok {
                rows, err := slot.drv.ListUsers(ctx); if err != nil { rows = []User{} }
                mu.Lock(); out[id] = rows; mu.Unlock()
            }
        }(id)
    }
    wg.Wait(); return out, nil
}

// InboundsAll: فقط گروه فعلی
func (g Group) InboundsAll(ctx context.Context) (map[string][]Inbound, error) {
    ids := g.IDs()
    out := make(map[string][]Inbound, len(ids))
    var mu sync.Mutex
    sem := make(chan struct{}, g.m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for _, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string) { defer wg.Done(); defer func(){<-sem}()
            if slot, ok := g.m.getSlot(id); ok {
                rows, err := slot.drv.ListInbounds(ctx); if err != nil { rows = []Inbound{} }
                mu.Lock(); out[id] = rows; mu.Unlock()
            }
        }(id)
    }
    wg.Wait(); return out, nil
}

// ExtMap: برگرداندن نقشه‌ی افزونه‌ها برای تمام پنل‌های گروه (فقط مواردی که پشتیبانی می‌کنند)
func (g Group) ExtMap[T any]() map[string]T {
    ids := g.IDs()
    out := make(map[string]T)
    for _, id := range ids { if v, ok := g.m.As[T](id); ok { out[id] = v } }
    return out
}

// TryEach: اجرای تابع روی تمام پنل‌های گروه که افزونه‌ی T را دارند
func (g Group) TryEach[T any](fn func(id string, t T) error) error {
    items := g.ExtMap[T]()
    var wg sync.WaitGroup
    sem := make(chan struct{}, g.m.opts.MaxConcurrency)
    errs := make(chan error, len(items))
    for id, v := range items {
        wg.Add(1); sem <- struct{}{}
        go func(id string, v T) { defer wg.Done(); defer func(){<-sem}()
            if err := fn(id, v); err != nil { errs <- err }
        }(id, v)
    }
    wg.Wait(); close(errs)
    for e := range errs { if e != nil { return e } }
    return nil
}

// میان‌بُرهای کاربردی برای خانواده‌ها
func (m *Manager) Marzban() Group    { return m.Kind(DriverMarzban) }
func (m *Manager) XUIAll() Group     { return m.Family("xui.") }
func (m *Manager) XUIGeneric() Group { return m.Kind(DriverXUIGeneric) }
func (m *Manager) XUISanaei() Group  { return m.Kind(DriverXUISanaei) }
func (m *Manager) XUIAlireza() Group { return m.Kind(DriverXUIAlireza) }


// HealthAll: اجرای health check روی تمام پنل‌های گروه (اگر پیاده‌سازی شده باشد)
func (g Group) HealthAll(ctx context.Context) map[string]error {
    ids := g.IDs()
    out := make(map[string]error, len(ids))
    var mu sync.Mutex
    sem := make(chan struct{}, g.m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for _, id := range ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string) {
            defer wg.Done(); defer func(){<-sem}()
            var err error
            if slot, ok := g.m.getSlot(id); ok {
                if h, ok2 := slot.drv.(HealthChecker); ok2 {
                    err = h.Health(ctx)
                }
            }
            mu.Lock(); out[id] = err; mu.Unlock()
        }(id)
    }
    wg.Wait(); return out
}
