package main

import (
	"fmt"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

func main() {
	mgr := core.New(core.WithRequestTimeout(6*time.Second))
	req := mgr.Request(core.WithReqTimeoutOpt(3 * time.Second))
	_, _ = req.UsersAll()
	fmt.Println("request ctx OK")
}
