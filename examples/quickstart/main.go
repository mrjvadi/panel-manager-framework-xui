package main

import (
    "context"
    "fmt"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
    _ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
)

func main() {
    mgr := core.New(core.WithRequestTimeout(8*time.Second))
    // placeholders — خودت مقادیر واقعی بذار
    _ = mgr.AttachByKind(core.PanelSpec{ ID:"xui-1", BaseURL:"https://xui.example.com", Auth: core.BasicAuth{ Username:"admin", Password:"pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverXUISanaei)
    _ = mgr.AttachByKind(core.PanelSpec{ ID:"marz-1", BaseURL:"https://marzban.example.com", Auth: core.BasicAuth{ Username:"admin", Password:"pass" }, TLS: core.TLS{ Insecure: true } }, core.DriverMarzban)

    // Panel session
    p := mgr.PanelSession("xui-1").WithTimeout(5*time.Second)
    inbs, _ := p.Inbounds()
    fmt.Println("inbounds (placeholder driver returns nil):", len(inbs))

    // Group filter
    ids := mgr.XUIAll().WhereVersionPrefix("v2.").IDs()
    fmt.Println("xui ids:", ids)

    // ctx usage example
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
    _ = ctx
}
