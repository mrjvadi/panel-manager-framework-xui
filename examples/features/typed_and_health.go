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

    _ = mgr.AttachByKind(core.PanelSpec{ ID: "marz-1", BaseURL: "https://marzban.example.com", Auth: core.BasicAuth{ Username: "admin", Password: "pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverMarzban)

    ctx := context.Background()
    // آماده‌سازی اتصال و کش نسخه
    _ = mgr.ConnectAll(ctx)

    // Typed API
    if mz, ok := mgr.ExtHub("marz-1").Marzban(); ok {
        if mzt, ok2 := mgr.As[ext.MarzbanTyped]("marz-1"); ok2 {
            sys, _ := mzt.SystemInfoTyped(ctx)
            fmt.Println("system version:", sys.Version)

            usages, _ := mzt.UsersUsageTyped(ctx)
            fmt.Println("typed users usage count:", len(usages))
        }
        // همچنان map-based هم در دسترس است
        _ , _ = mz.UsersUsage(ctx)
    }

    // Health check گروهی
    health := mgr.Marzban().HealthAll(ctx)
    for id, err := range health {
        if err != nil { fmt.Println("panel", id, "UNHEALTHY:", err) } else { fmt.Println("panel", id, "OK") }
    }

    // فیلتر براساس نسخهٔ کش‌شده از system
    v2 := mgr.Marzban().WhereVersionPrefix("v2.")
    fmt.Println("marzban v2.*:", v2.IDs())
}
