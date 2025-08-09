package main

import (
    "fmt"
    "time"
    "context"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
    mgr := core.New(core.WithBaseContext(context.Background()), core.WithRequestTimeout(8*time.Second))
    _ = mgr.AttachByKind(core.PanelSpec{ ID: "xui-s-1", BaseURL: "https://xui.example.com", Auth: core.BasicAuth{ Username: "admin", Password: "pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverXUISanaei)

    p := mgr.PanelSession("xui-s-1")

    // بدون ctx: کلون با پورت مشخص
    inb, _ := p.XUI().CloneInboundWithPort(1, 24443)
    fmt.Println("cloned:", inb.ID, inb.Port)

    // بدون ctx: کلون با اسم مشخص
    inb2, _ := p.XUI().CloneInbound(1, xdto.CloneInboundOptions{ Remark: ptrS("copy-special") })
    fmt.Println("cloned2:", inb2.Remark, inb2.Port)
}

func ptrS(s string) *string { return &s }
