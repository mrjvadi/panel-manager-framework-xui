package main

import (
    "fmt"
    "time"
    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
    mgr := core.New(core.WithRequestTimeout(8*time.Second))
    _ = mgr.AttachByKind(core.PanelSpec{ ID:"xui-1", BaseURL:"https://xui.example.com", Auth: core.BasicAuth{ Username:"admin", Password:"pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverXUISanaei)
    p := mgr.PanelSession("xui-1").WithTimeout(5*time.Second)
    inb, _ := p.XUI().CloneInbound(1, xdto.CloneInboundOptions{})
    fmt.Println("cloned inbound:", inb.ID, inb.Port)
}
