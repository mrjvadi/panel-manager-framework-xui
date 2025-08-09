package core

import "strings"
import "sync"
import "context"

type Group struct { m *Manager; ids []string }

func (m *Manager) allIDs() []string {
    m.mu.RLock(); defer m.mu.RUnlock()
    out := make([]string,0,len(m.panels))
    for id, s := range m.panels { if s.enabled { out = append(out, id) } }
    return out
}

func (m *Manager) XUIAll() Group {
    ids := m.allIDs()
    out := make([]string,0,len(ids))
    for _, id := range ids {
        if s, ok := m.getSlot(id); ok { if strings.HasPrefix(s.drv.Name(), "xui.") { out = append(out, id) } }
    }
    return Group{ m:m, ids: out }
}

func (m *Manager) VersionEq(v string) Group {
    ids := m.allIDs()
    out := make([]string,0,len(ids))
    for _, id := range ids {
        if s, ok := m.getSlot(id); ok { if compareVersionStr(s.drv.Version(), v) == 0 { out = append(out, id) } }
    }
    return Group{ m:m, ids: out }
}

func (g Group) IDs() []string { return append([]string(nil), g.ids...) }

func (g Group) WhereVersionPrefix(p string) Group {
    out := make([]string,0,len(g.ids))
    for _, id := range g.ids {
        if s, ok := g.m.getSlot(id); ok { if strings.HasPrefix(s.drv.Version(), p) { out = append(out, id) } }
    }
    return Group{ m:g.m, ids: out }
}

func (g Group) UsersAll(ctx context.Context) map[string][]User {
    out := map[string][]User{}
    var mu sync.Mutex
    sem := make(chan struct{}, g.m.opts.MaxConcurrency)
    var wg sync.WaitGroup
    for _, id := range g.ids {
        wg.Add(1); sem <- struct{}{}
        go func(id string){
            defer wg.Done(); defer func(){<-sem}()
            rows, _ := g.m.Users(ctx, id)
            mu.Lock(); out[id]=rows; mu.Unlock()
        }(id)
    }
    wg.Wait(); return out
}
