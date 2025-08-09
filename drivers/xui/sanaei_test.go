package xui

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
)

func TestCloneInboundTyped(t *testing.T) {
    var nextID = 100
    var token = "X"
    attempts := 0

    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        switch {
        case r.URL.Path == "/login" && r.Method == "POST":
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(`{"token":"`+token+`"}`))
        case r.URL.Path == "/panel/api/inbounds/get/1":
            if r.Header.Get("Authorization") != "Bearer "+token { w.WriteHeader(401); return }
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(`{"id":1,"remark":"base","protocol":"vless","port":443,"settings":{},"sniffing":{},"streamSettings":{}}`))
        case r.URL.Path == "/panel/api/inbounds/add" && r.Method == "POST":
            if r.Header.Get("Authorization") != "Bearer "+token { w.WriteHeader(401); return }
            attempts++
            if attempts == 1 {
                // simulate conflict to trigger retry
                w.WriteHeader(409); return
            }
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(fmt.Sprintf(`{"id":%d,"remark":"ok","protocol":"vless","port":24443}`, nextID)))
        default:
            w.WriteHeader(404)
        }
    }))
    defer srv.Close()

    sp := core.PanelSpec{ ID:"x1", BaseURL: srv.URL, Auth: core.BasicAuth{ Username:"a", Password:"b" }, TLS: core.TLS{ Insecure: true } }
    drv, err := NewSanaei(sp, core.WithHTTPClient(srv.Client()), core.WithRequestTimeout(5*time.Second))
    if err != nil { t.Fatal(err) }

    xt := drv.(ext.XUITyped)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()

    inb, err := xt.CloneInboundTyped(ctx, 1, xdto.CloneInboundOptions{ /* none: random */ })
    if err != nil { t.Fatal(err) }
    if inb.ID == 0 || inb.Port == 0 { t.Fatalf("invalid cloned inbound: %+v", inb) }
}
