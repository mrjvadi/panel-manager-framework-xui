package core

import (

    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// HTTPError: خطای سطح HTTP با کد و بدنه‌ی کوتاه
type HTTPError struct {
    Code int
    Body string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("http error: status=%d body=%s", e.Code, e.Body)
}

func IsHTTPStatus(err error, code int) bool {
    if he, ok := err.(*HTTPError); ok {
        return he.Code == code
    }
    return false
}

// DoJSON: ارسال درخواست و دیکد پاسخ JSON با retry + breaker
func DoJSON(ctx context.Context, cli *HTTP, req *http.Request, v any) error {
    if cli == nil || cli.Client == nil {
        return fmt.Errorf("nil http client")
    }
    if req == nil {
        return fmt.Errorf("nil request")
    }
    attempts := cli.Retry.MaxAttempts
    if attempts <= 0 { attempts = 1 }

    for i := 1; i <= attempts; i++ {
        if cli.Br != nil && !cli.Br.Allow() {
            return ErrCircuitOpen
        }
        resp, err := cli.Client.Do(req.Clone(ctx))
        if err != nil {
            if cli.Br != nil { cli.Br.OnFailure() }
            if i == attempts || !retryableNetErr(err) { return fmt.Errorf("request %s %s failed: %w", req.Method, req.URL, err) }
            time.Sleep(cli.Retry.backoff(i))
            continue
        }
        func() {
            defer resp.Body.Close()
            if resp.StatusCode < 200 || resp.StatusCode >= 300 {
                b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
                he := &HTTPError{ Code: resp.StatusCode, Body: string(b) }
                if cli.Br != nil { cli.Br.OnFailure() }
                if i == attempts || !cli.Retry.Codes[resp.StatusCode] {
                    err = he
                    return
                }
                // retryable http status
                err = he
            } else {
                if v == nil {
                    io.Copy(io.Discard, resp.Body)
                } else {
                    if de := json.NewDecoder(resp.Body).Decode(v); de != nil {
                        if cli.Br != nil { cli.Br.OnFailure() }
                        err = fmt.Errorf("decode json: %w", de)
                        return
                    }
                }
                if cli.Br != nil { cli.Br.OnSuccess() }
                err = nil
            }
        }()
        if err == nil {
            return nil
        }
        if _, ok := err.(*HTTPError); ok {
            if i == attempts || !cli.Retry.Codes[err.(*HTTPError).Code] {
                return err
            }
            time.Sleep(cli.Retry.backoff(i))
            continue
        }
        return err
    }
    return fmt.Errorf("unreachable DoJSON loop")
}
