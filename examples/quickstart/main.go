package main

import (
	"fmt"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func main() {
	mgr := core.New(core.WithRequestTimeout(8 * time.Second))
	fmt.Println("manager ready:", mgr != nil)
}
