package integration

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func randPort() int { return 20000 + rand.Intn(40000) }

func Test_XUI_Sanaei_CloneInbound(t *testing.T) {
	base := "http://91.107.190.148:37736/C8TbUbb7qLxNDEzrGs"
	user := "4R9q2FXp0i"
	pass := "ntrdglyA8J"
	baseInboundID := "1"
	if base == "" || user == "" || pass == "" || baseInboundID == "" {
		t.Skip("set PMF_XUI_SANAEI_URL/USER/PASS and PMF_XUI_BASE_INBOUND_ID")
	}
	inboundID, err := strconv.Atoi(baseInboundID)
	if err != nil || inboundID <= 0 {
		t.Fatalf("invalid PMF_XUI_BASE_INBOUND_ID: %v", baseInboundID)
	}

	mgr := core.New(core.WithTimeout(20*time.Second), core.WithMaxConcurrency(16))

	spec := core.PanelSpec{
		ID:      "xui-sanaei-live",
		BaseURL: base,
		Auth:    core.BasicAuth{Username: user, Password: pass},
		TLS:     core.TLS{Insecure: true},
	}
	if err := mgr.AttachByKind(spec, core.DriverXUISanaei); err != nil {
		t.Fatalf("attach: %v", err)
	}

	xt, ok := core.As[ext.XUITyped](mgr, "xui-sanaei-live")
	if !ok {
		t.Fatal("XUITyped not supported")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout
	defer cancel()

	p := mgr.PanelSession("xui-sanaei-live")

	info, err := xt.GetInboundTyped(ctx, inboundID)
	if err != nil {
		t.Fatalf("list inbounds failed: %v", err)
	}
	t.Logf("found %v inbounds", info)

	// A light update
	remark := "my-remark2"
	port := 23256
	client := xdto.ClientCreate{
		Email:      "user2@example.com",
		Enable:     true,
		TotalGB:    50, // GB
		ExpiryTime: 0,  // Unix seconds or 0
		LimitIP:    0,
		Comment:    "created by tests",
	}

	cloneda, err := p.XUI().CloneInboundShallow(1, xdto.CloneInboundOptions{
		Remark: &remark,
		Port:   &port,
		Client: &client, 
	})
	if err != nil {
		t.Fatalf("clone inbound shallow failed: %v", err)

	}
	t.Logf("updated inbound id=%d remark=%s", cloneda.ID, cloneda.Remark)

}