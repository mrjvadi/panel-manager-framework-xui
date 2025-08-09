//go:build examples
// +build examples

package main

import (
    "fmt"
    "time"
    "context"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
)

func main() {
    mgr := core.New(
        core.WithMaxConcurrency(16),
        core.WithBaseContext(context.Background()),       // یکبار تنظیم می‌کنی
        core.WithRequestTimeout(10*time.Second),          // تایم‌اوت پیش‌فرض هر درخواست
    )

    _ = mgr.AttachByKind(core.PanelSpec{ ID: "marz-1", BaseURL: "https://marzban.example.com", Auth: core.BasicAuth{ Username: "admin", Password: "pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverMarzban)

    // اسکوپ درخواست بساز (فقط یک بار)
    r := mgr.Request(core.WithReqTimeoutOpt(8*time.Second)).WithValue("tenant","acme")

    // حالا بدون ساخت ctx در هر بار، فراخوانی کن
    usersByPanel, _ := r.UsersAll()
    fmt.Println("panels:", len(usersByPanel))

    // اجرای اکستنشن‌ها با TryEachCtx و ctx تزریق‌شده
    _ = r.Marzban().TryEachCtx[ext.Marzban](func(ctx context.Context, id string, mz ext.Marzban) error {
        _, _ = mz.UsersUsage(ctx) // ctx آماده است
        fmt.Println("ok:", id)
        return nil
    })
}
