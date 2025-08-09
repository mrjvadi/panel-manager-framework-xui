package core

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type HTTPError struct{ Code int; Body string }
func (e *HTTPError) Error() string { return fmt.Sprintf("http error: %d %s", e.Code, e.Body) }
func IsHTTPStatus(err error, code int) bool { if he, ok := err.(*HTTPError); ok { return he.Code==code }; return false }

func DoJSON(ctx context.Context, cli *HTTP, req *http.Request, v any) error {
    if cli == nil || cli.Client == nil { return fmt.Errorf("nil client") }
    attempts := cli.Retry.MaxAttempts; if attempts <= 0 { attempts = 1 }
    for i := 1; i <= attempts; i++ {
        if cli.Br != nil && !cli.Br.Allow() { return ErrCircuitOpen }
        if cli.Log != nil { cli.Log.Debug("http.attempt", "n", i, "method", req.Method, "url", req.URL.String()) }
        resp, err := cli.Client.Do(req.Clone(ctx))
        if err != nil {
            if cli.Log != nil { cli.Log.Warn("http.error", "err", err) }
            if i == attempts || !retryableNetErr(err) { return fmt.Errorf("request failed: %w", err) }
            time.Sleep(cli.Retry.backoff(i)); continue
        }
        var retErr error
        func() {
            defer resp.Body.Close()
            if resp.StatusCode < 200 || resp.StatusCode >= 300 {
                b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
                he := &HTTPError{ Code: resp.StatusCode, Body: string(b) }
                if cli.Log != nil { cli.Log.Warn("http.status", "code", resp.StatusCode) }
                if i == attempts || !cli.Retry.Codes[resp.StatusCode] { retErr = he; return }
                retErr = he; return
            }
            if v == nil { io.Copy(io.Discard, resp.Body) } else {
                if de := json.NewDecoder(resp.Body).Decode(v); de != nil { retErr = fmt.Errorf("decode json: %w", de); return }
            }
            if cli.Log != nil { cli.Log.Info("http.ok", "code", resp.StatusCode) }
            retErr = nil
        }()
        if retErr == nil { return nil }
        if _, ok := retErr.(*HTTPError); ok { time.Sleep(cli.Retry.backoff(i)); continue }
        return retErr
    }
    return fmt.Errorf("unreachable")
}
