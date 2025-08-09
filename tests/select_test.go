package tests

import (
    "testing"
    "context"
    core "github.com/mrjvadi/panel-manager-framework-xui/core"
)

type dummyDriver struct{ name, ver string; cap core.Capabilities }
func (d dummyDriver) Name() string { return d.name }
func (d dummyDriver) Version() string { return d.ver }
func (d dummyDriver) Capabilities() core.Capabilities { return d.cap }
func (d dummyDriver) Login(ctx context.Context) error { return nil }
func (d dummyDriver) ListUsers(ctx context.Context) ([]core.User, error) { return nil, nil }
func (d dummyDriver) ListInbounds(ctx context.Context) ([]core.Inbound, error) { return nil, nil }

func makeFactory(name, ver string) core.FactoryFn {
    return func(sp core.PanelSpec, _ ...core.Option) (core.Driver, error) {
        return dummyDriver{name:name, ver:ver}, nil
    }
}

func TestVersionFilters(t *testing.T) {
    core.Register("dummy.a", makeFactory("xui.sanaei","v2.6.1"))
    core.Register("dummy.b", makeFactory("xui.alireza","v2.2.0"))
    core.Register("dummy.c", makeFactory("marzban","v0.8.4"))

    m := core.New()
    _ = m.Attach(core.PanelSpec{ID:"a"}, "dummy.a")
    _ = m.Attach(core.PanelSpec{ID:"b"}, "dummy.b")
    _ = m.Attach(core.PanelSpec{ID:"c"}, "dummy.c")

    got := m.XUIAll().WhereVersionPrefix("v2.").IDs()
    if len(got) != 2 { t.Fatalf("expected 2 xui, got %v", got) }

    got2 := m.VersionEq("v0.8.4").IDs()
    if len(got2) != 1 || got2[0] != "c" { t.Fatalf("want c, got %v", got2) }
}
