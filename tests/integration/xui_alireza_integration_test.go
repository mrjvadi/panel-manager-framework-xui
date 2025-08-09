package integration

//
//import (
//	"context"
//	"math/rand"
//	"net/http"
//	"strconv"
//	"testing"
//	"time"
//
//	"github.com/mrjvadi/panel-manager-framework-xui/core"
//	xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
//	ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
//	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
//)
//
//func Test_XUI_Alireza_CloneInbound(t *testing.T) {
//	base := "http://91.107.190.148:6522/P7HKO0xVHMfklTv"
//	user := "admin"
//	pass := "admin"
//	baseInboundID := "1"
//	if base == "" || user == "" || pass == "" || baseInboundID == "" {
//		t.Skip("set PMF_XUI_ALIREZA_URL/USER/PASS and PMF_XUI_BASE_INBOUND_ID")
//	}
//	inboundID, err := strconv.Atoi(baseInboundID)
//	if err != nil || inboundID <= 0 {
//		t.Fatalf("invalid PMF_XUI_BASE_INBOUND_ID: %v", baseInboundID)
//	}
//
//	mgr := core.New(core.WithTimeout(20*time.Second), core.WithMaxConcurrency(16))
//
//	spec := core.PanelSpec{
//		ID:      "xui-alireza-live",
//		BaseURL: base,
//		Auth:    core.BasicAuth{Username: user, Password: pass},
//		TLS:     core.TLS{Insecure: true},
//	}
//	if err := mgr.AttachByKind(spec, core.DriverXUIAlireza); err != nil {
//		t.Fatalf("attach: %v", err)
//	}
//
//	xt, ok := core.As[ext.XUITyped](mgr, "xui-alireza-live")
//	if !ok {
//		t.Fatal("XUITyped not supported")
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	// ریتری ساده روی 409
//	try := 0
//	var cloned xdto.Inbound
//	for {
//		try++
//		p := 20000 + rand.Intn(40000)
//		out, err := xt.CloneInboundTyped(ctx, inboundID, xdto.CloneInboundOptions{Port: &p})
//		if err == nil {
//			cloned = out
//			break
//		}
//		if he, ok := err.(*core.HTTPError); ok && he.Code == http.StatusConflict && try < 5 {
//			t.Logf("conflict on port %d, retrying...", p)
//			continue
//		}
//		t.Fatalf("clone failed: %v", err)
//	}
//	t.Logf("cloned inbound id=%d port=%d", cloned.ID, cloned.Port)
//}
