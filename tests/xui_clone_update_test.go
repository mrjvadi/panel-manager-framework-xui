package tests

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    core "github.com/mrjvadi/panel-manager-framework-xui/core"
    xdto "github.com/mrjvadi/panel-manager-framework-xui/core/dto/xui"
    ext "github.com/mrjvadi/panel-manager-framework-xui/core/ext"
    xui "github.com/mrjvadi/panel-manager-framework-xui/drivers/xui"
)

func TestXUI_CloneAndUpdateInboundTyped(t *testing.T) {
    var token = "X"
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
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(fmt.Sprintf(`{"id":%d,"remark":"ok","protocol":"vless","port":24443}`, 100)))
        case r.URL.Path == "/panel/api/inbounds/update/100" && r.Method == "POST":
            if r.Header.Get("Authorization") != "Bearer "+token { w.WriteHeader(401); return }
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(`{"id":100,"remark":"changed","protocol":"vless","port":24444}`))
        default:
            w.WriteHeader(404)
        }
    }))
    defer srv.Close()

    sp := core.PanelSpec{ ID:"x1", BaseURL: srv.URL, Auth: core.BasicAuth{ Username:"a", Password:"b" }, TLS: core.TLS{ Insecure: true } }
    drv, err := xui.NewSanaei(sp, core.WithHTTPClient(srv.Client()), core.WithRequestTimeout(5*time.Second))
    if err != nil { t.Fatal(err) }
    xt := drv.(ext.XUITyped)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()

    cloned, err := xt.CloneInboundTyped(ctx, 1, xdto.CloneInboundOptions{})
    if err != nil { t.Fatal(err) }
    if cloned.ID == 0 { t.Fatalf("invalid cloned inbound: %+v", cloned) }

    updated, err := xt.UpdateInboundTyped(ctx, cloned.ID, xdto.InboundUpdate{ Remark: "changed", Protocol: "vless", Port: 24444 })
    if err != nil { t.Fatal(err) }
    if updated.ID != 100 { t.Fatalf("update not applied: %+v", updated) }
}
