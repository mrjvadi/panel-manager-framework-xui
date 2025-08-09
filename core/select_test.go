package core

import "testing"

type dummyDriver struct{ name, ver string; cap Capabilities }
func (d dummyDriver) Name() string { return d.name }
func (d dummyDriver) Version() string { return d.ver }
func (d dummyDriver) Capabilities() Capabilities { return d.cap }
func (d dummyDriver) Login(ctx context.Context) error { return nil }
func (d dummyDriver) ListUsers(ctx context.Context) ([]User, error) { return nil, nil }
func (d dummyDriver) ListInbounds(ctx context.Context) ([]Inbound, error) { return nil, nil }
func (d dummyDriver) CreateUser(ctx context.Context, u User) (User, error) { return u, ErrNotImplemented }
func (d dummyDriver) UpdateUser(ctx context.Context, u User) (User, error) { return u, ErrNotImplemented }
func (d dummyDriver) DeleteUser(ctx context.Context, id string) error { return ErrNotImplemented }
func (d dummyDriver) SuspendUser(ctx context.Context, id string) error { return ErrNotImplemented }
func (d dummyDriver) ResumeUser(ctx context.Context, id string) error { return ErrNotImplemented }
func (d dummyDriver) ResetUserTraffic(ctx context.Context, id string) error { return ErrNotImplemented }
func (d dummyDriver) CreateInbound(ctx context.Context, in Inbound) (Inbound, error) { return in, ErrNotImplemented }
func (d dummyDriver) UpdateInbound(ctx context.Context, in Inbound) (Inbound, error) { return in, ErrNotImplemented }
func (d dummyDriver) DeleteInbound(ctx context.Context, id string) error { return ErrNotImplemented }

import "context"

func TestVersionFilters(t *testing.T) {
    m := New()
    m.mu.Lock()
    m.panels["a"] = &panelSlot{ drv: dummyDriver{name:"xui.sanaei", ver:"v2.6.1"}, enabled: true }
    m.panels["b"] = &panelSlot{ drv: dummyDriver{name:"xui.alireza", ver:"v2.2.0"}, enabled: true }
    m.panels["c"] = &panelSlot{ drv: dummyDriver{name:"marzban", ver:"v0.8.4"}, enabled: true }
    m.mu.Unlock()

    got := m.XUIAll().WhereVersionPrefix("v2.").IDs()
    if len(got) != 2 { t.Fatalf("expected 2 xui, got %v", got) }

    got2 := m.VersionEq("v0.8.4").IDs()
    if len(got2) != 1 || got2[0] != "c" { t.Fatalf("want c, got %v", got2) }
}
