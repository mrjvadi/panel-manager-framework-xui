package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	core "github.com/mrjvadi/panel-manager-framework-xui/core"
	driver "github.com/mrjvadi/panel-manager-framework-xui/drivers/marzban"
)

type hasSystem interface {
	SystemInfo(ctx context.Context) (map[string]any, error)
}

func TestMarzban_ListUsers_And_SystemVersion(t *testing.T) {
	token := "T1"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/admin/token":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"access_token":"` + token + `"}`))
		case "/api/admin/users":
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(401)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"users":[{"username":"ali","up":1,"down":2}]}`))
		case "/api/system":
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(401)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"version":"v0.9.0","uptime":123}`))
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	sp := core.PanelSpec{ID: "m1", BaseURL: srv.URL, Auth: core.BasicAuth{Username: "u", Password: "p"}, TLS: core.TLS{Insecure: true}}
	d, err := driver.New(sp, core.WithRequestTimeout(5*time.Second), core.WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatal(err)
	}

	arr, err := d.ListUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(arr) != 1 || arr[0].Username != "ali" {
		t.Fatalf("bad users: %#v", arr)
	}

	ms, ok := d.(hasSystem)
	if !ok {
		t.Fatal("driver does not expose SystemInfo")
	}

	if _, err := ms.SystemInfo(context.Background()); err != nil {
		t.Fatal(err)
	}
	if d.Version() != "v0.9.0" {
		t.Fatalf("runtime version not cached, got %s", d.Version())
	}
}
