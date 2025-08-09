package main

import (
    "context"
    "fmt"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
    mgr := core.New(core.WithTimeout(20*time.Second), core.WithMaxConcurrency(16))

    _ = mgr.AttachByKind(core.PanelSpec{
        ID: "marz-1",
        BaseURL: "https://marzban.example.com",
        Auth: core.BasicAuth{ Username: "admin", Password: "pass" },
        TLS: core.TLS{ Insecure: true },
    }, core.DriverMarzban)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    usersAll, _ := mgr.UsersAll(ctx)
    for id, users := range usersAll { fmt.Println("panel:", id, "users:", len(users)) }

    go func() {
        time.Sleep(2 * time.Second)
        _ = mgr.AttachByKind(core.PanelSpec{
            ID: "xui-s-1",
            BaseURL: "https://xui-sanaei.example.com",
            Auth: core.BasicAuth{ Username: "admin", Password: "pass" },
            TLS: core.TLS{ Insecure: true },
            Version: "v2.6.1",
        }, core.DriverXUISanaei)
        fmt.Println("attached xui-s-1 at runtime")
    }()

    go func() {
        time.Sleep(3 * time.Second)
        mgr.Disable("marz-1"); fmt.Println("disabled marz-1")
        time.Sleep(2 * time.Second)
        mgr.Enable("marz-1"); fmt.Println("enabled marz-1")
    }()

    time.Sleep(6 * time.Second)

    usersAll, _ = mgr.UsersAll(ctx)
    for id, users := range usersAll { fmt.Println("panel:", id, "users:", len(users)) }
}
