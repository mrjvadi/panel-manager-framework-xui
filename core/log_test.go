package core

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

type memLogger struct{ count int }
func (m *memLogger) Debug(string, ...any) { m.count++ }
func (m *memLogger) Info(string, ...any)  { m.count++ }
func (m *memLogger) Warn(string, ...any)  { m.count++ }
func (m *memLogger) Error(string, ...any) { m.count++ }

func TestDoJSONLogging(t *testing.T) {
    ml := &memLogger{}
    cli := NewHTTP("http://x", true, 2*time.Second, nil)
    cli.Log = ml
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type","application/json")
        json.NewEncoder(w).Encode(map[string]any{"ok":true})
    }))
    defer srv.Close()
    req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
    var out map[string]any
    if err := DoJSON(context.Background(), cli, req, &out); err != nil { t.Fatal(err) }
    if ml.count == 0 { t.Fatal("expected some logs") }
}
