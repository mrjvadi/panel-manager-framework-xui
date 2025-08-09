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

func TestXUI_CloneInboundTyped(t *testing.T) {
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
            if attempts == 1 { w.WriteHeader(409); return }
            w.Header().Set("Content-Type","application/json")
            w.Write([]byte(fmt.Sprintf(`{"id":%d,"remark":"ok","protocol":"vless","port":24443}`, 100)))
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
    inb, err := xt.CloneInboundTyped(ctx, 1, xdto.CloneInboundOptions{})
    if err != nil { t.Fatal(err) }
    if inb.ID == 0 { t.Fatalf("invalid cloned inbound: %+v", inb) }
}
