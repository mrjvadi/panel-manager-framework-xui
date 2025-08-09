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
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
    mgr := core.New(core.WithMaxConcurrency(16))
    _ = mgr.AttachByKind(core.PanelSpec{ ID: "xui-s-1", BaseURL: "https://xui-sanaei.example.com", Auth: core.BasicAuth{ Username: "admin", Password: "pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverXUISanaei)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if xt, ok := mgr.As[ext.XUITyped]("xui-s-1"); ok {
        inb, _ := xt.GetInboundTyped(ctx, 1)
        fmt.Println("Inbound:", inb.ID, inb.Protocol, inb.Port)

        // create/update example (payload fake)
        _, _ = xt.AddInboundTyped(ctx, xdto.InboundCreate{ Protocol: "vless", Port: 443 })
    }
}
