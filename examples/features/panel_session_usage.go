//go:build examples
// +build examples

package main

import (
    "fmt"
    "time"
    "context"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    mdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/marzban"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
)

func main() {
    mgr := core.New(
        core.WithBaseContext(context.Background()),
        core.WithRequestTimeout(8*time.Second),
    )

    _ = mgr.AttachByKind(core.PanelSpec{
        ID: "marz-1",
        BaseURL: "https://marzban.example.com",
        Auth: core.BasicAuth{ Username: "admin", Password: "pass" },
        TLS: core.TLS{ Insecure: true },
    }, core.DriverMarzban)

    // یک بار پنل را باز کن
    p := mgr.PanelSession("marz-1").WithTimeout(5*time.Second)

    // بدون ساخت ctx
    users, _ := p.Users()
    fmt.Println("users:", len(users))

    // اکستنشن Marzban (typed اگر موجود بود)
    _ = p.Try[ext.Marzban](func(ctx context.Context, mz ext.Marzban) error {
        usg, _ := mz.UsersUsage(ctx)
        // تبدیل map->struct با Binder
        if len(usg) > 0 {
            type Row struct { Username string `json:"username"`; Up int64 `json:"up"`; Down int64 `json:"down"` }
            r, _ := core.MapInto[Row](usg[0])
            fmt.Println("first:", r.Username, r.Up+r.Down)
        }
        return nil
    })

    // اگر MarzbanTyped در دسترس بود، بدون Binder
    if mzt, ok := p.As[ext.MarzbanTyped](); ok {
        arr, _ := mzt.UsersUsageTyped(context.Background())
        fmt.Println("typed count:", len(arr))
        // یا با ctx داخلی:
        _ = p.Try[ext.MarzbanTyped](func(ctx context.Context, t ext.MarzbanTyped) error {
            sys, _ := t.SystemInfoTyped(ctx)
            fmt.Println("version:", sys.Version)
            return nil
        })
    }

    // مثال Binder مستقیم
    src := map[string]any{"username":"ali","up":10,"down":20}
    var dst mdto.UserUsage
    _ = core.Bind().From(src).Into(&dst)
    fmt.Println("dto total:", dst.Up+dst.Down)
}
