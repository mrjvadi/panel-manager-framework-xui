package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
	mgr := core.New(core.WithRequestTimeout(8*time.Second))
	_ = mgr.AttachByKind(core.PanelSpec{
		ID:      "xui-1",
		BaseURL: "https://xui.example",
		Auth:    core.BasicAuth{Username: "u", Password: "p"},
		TLS:     core.TLS{Insecure: true},
	}, core.DriverXUISanaei)

	if xt, ok := core.As[ext.XUITyped](mgr, "xui-1"); ok {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = xt.GetInboundTyped(ctx, 1)
	}
	fmt.Println("xui typed OK")
}
