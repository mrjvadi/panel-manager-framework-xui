package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
	mgr := core.New(core.WithTimeout(20*time.Second), core.WithMaxConcurrency(16))

	_ = mgr.AttachByKind(core.PanelSpec{ID: "marz-1", BaseURL: "https://marzban.example.com", Auth: core.BasicAuth{Username: "admin", Password: "pass"}, TLS: core.TLS{Insecure: true}}, core.DriverMarzban)
	_ = mgr.AttachByKind(core.PanelSpec{ID: "marz-2", BaseURL: "https://m2.example.com", Auth: core.BasicAuth{Username: "a", Password: "b"}, TLS: core.TLS{Insecure: true}}, core.DriverMarzban)
	_ = mgr.AttachByKind(core.PanelSpec{ID: "xui-s-1", BaseURL: "https://xui-sanaei.example.com", Auth: core.BasicAuth{Username: "x", Password: "y"}, TLS: core.TLS{Insecure: true}, Version: "v2.6.1"}, core.DriverXUISanaei)
	_ = mgr.AttachByKind(core.PanelSpec{ID: "xui-a-1", BaseURL: "https://xui-alireza.example.com", Auth: core.BasicAuth{Username: "x", Password: "y"}, TLS: core.TLS{Insecure: true}, Version: "v2.3.0"}, core.DriverXUIAlireza)

	ctx := context.Background()

	// --- دسترسی تک-پنلی ---
	if mz, ok := mgr.ExtHub("marz-1").Marzban(); ok {
		u, _ := mz.UserUsage(ctx, "ali")
		fmt.Println("ExtHub().Marzban() => ali:", u)
	}

	// --- گروهی بر اساس نوع پنل ---
	g := mgr.Marzban()
	fmt.Println("Marzban IDs:", g.IDs())

	usersByPanel, _ := g.UsersAll(ctx)
	for id, us := range usersByPanel {
		fmt.Println("panel:", id, "users:", len(us))
	}

	_ = g.TryEach[ext.Marzban](func(id string, mz ext.Marzban) error {
		m, _ := mz.UsersUsage(ctx)
		fmt.Println("panel:", id, "usersUsage items:", len(m))
		return nil
	})

	// --- خانواده‌ی XUI جدا ---
	xAll := mgr.XUIAll()
	xS := mgr.XUISanaei()
	xA := mgr.XUIAlireza()
	xG := mgr.XUIGeneric()

	fmt.Println("XUI All:", xAll.IDs())
	fmt.Println("XUI Sanaei:", xS.IDs())
	fmt.Println("XUI Alireza:", xA.IDs())
	fmt.Println("XUI Generic:", xG.IDs())

	_ = xS.TryEach[ext.XUI](func(id string, x ext.XUI) error {
		ips, _ := x.ClientIPs(ctx, "user@example.com")
		fmt.Println("panel:", id, "ips:", len(ips))
		return nil
	})

	// --- فیلتر ورژنی ---
	only261 := mgr.VersionEq("v2.6.1")
	fmt.Println("Version v2.6.1:", only261.IDs())

	xFamilyV2 := mgr.XUIAll().WhereVersionPrefix("v2.")
	fmt.Println("XUI v2.*:", xFamilyV2.IDs())
}
