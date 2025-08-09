package integration

import (
	"context"
	"github.com/mrjvadi/panel-manager-framework-xui/core"
	_ "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
	"testing"
	"time"
)

type hasSystem interface {
	SystemInfo(ctx context.Context) (map[string]any, error)
}

func Test_Marzban_System_And_ListUsers(t *testing.T) {
	base := "https://bots.saitsazs.ir:2053"
	user := "Et_MrJavdi"
	pass := "Et_MrJavdi"
	if base == "" || user == "" || pass == "" {
		t.Skip("set PMF_MARZBAN_URL/USER/PASS to run this integration test")
	}

	insecure := true

	mgr := core.New(core.WithTimeout(20*time.Second), core.WithMaxConcurrency(16))

	_ = mgr.AttachByKind(core.PanelSpec{
		ID:      "marz-live",
		BaseURL: base,
		Auth:    core.BasicAuth{Username: user, Password: pass},
		TLS:     core.TLS{Insecure: insecure},
	}, core.DriverMarzban)

	ctx := context.Background()

	users, err := mgr.Users(ctx, "marz-live")
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	t.Logf("users: %d", len(users))

	// Optional: اگر درایور SystemInfo را اکسپوز کرده باشد
	if d, ok := core.As[hasSystem](mgr, "marz-live"); ok {
		sys, err := d.SystemInfo(ctx)
		if err != nil {
			t.Fatalf("system info: %v", err)
		}
		t.Logf("system: %#v", sys)
	}
}
