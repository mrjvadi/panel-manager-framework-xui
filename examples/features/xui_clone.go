//go:build examples
// +build examples

package main

import (
    "context"
    "fmt"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
    mgr := core.New(core.WithRequestTimeout(10*time.Second))
    _ = mgr.AttachByKind(core.PanelSpec{ ID: "xui-s-1", BaseURL: "https://xui.example.com", Auth: core.BasicAuth{ Username: "admin", Password: "pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverXUISanaei)

    if xt, ok := mgr.As[ext.XUITyped]("xui-s-1"); ok {
        ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
        defer cancel()

        // مثال: کلون با پورت مشخص
        p := 20543
        inb, err := xt.CloneInboundTyped(ctx, 1, xdto.CloneInboundOptions{ Port: &p })
        fmt.Println("cloned:", inb.ID, inb.Port, "err:", err)

        // مثال: کلون با اسم مشخص
        name := "my-inb-copy"
        inb2, _ := xt.CloneInboundTyped(ctx, 1, xdto.CloneInboundOptions{ Remark: &name })
        fmt.Println("cloned2:", inb2.Remark, inb2.Port)
    }
}
