package main

import (
	"fmt"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
)

func main() {
	mgr := core.New(core.WithRequestTimeout(10*time.Second))
	_ = mgr.AttachByKind(core.PanelSpec{
		ID:      "mz-1",
		BaseURL: "https://marzban.example",
		Auth:    core.BasicAuth{Username: "admin", Password: "pass"},
	}, core.DriverMarzban)

	p := mgr.PanelSession("mz-1").WithTimeout(5 * time.Second)
	_, _ = p.Users()
	fmt.Println("panel session OK")
}
