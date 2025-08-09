package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPError struct {
	Code int
	Body string
}

func (e *HTTPError) Error() string { return fmt.Sprintf("http error: %d %s", e.Code, e.Body) }
func IsHTTPStatus(err error, code int) bool {
	if he, ok := err.(*HTTPError); ok {
		return he.Code == code
	}
	return false
}

func DoJSON(ctx context.Context, cli *HTTP, req *http.Request, v any) error {
	if cli == nil || cli.Client == nil {
		return fmt.Errorf("nil client")
	}
	attempts := cli.Retry.MaxAttempts
	if attempts <= 0 {
		attempts = 1
	}

	var lastErr error

	for i := 1; i <= attempts; i++ {
		if cli.Br != nil && !cli.Br.Allow() {
			return ErrCircuitOpen
		}
		if cli.Log != nil {
			cli.Log.Debug("http.attempt", "n", i, "method", req.Method, "url", req.URL.String())
		}

		resp, err := cli.Client.Do(req.Clone(ctx))
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if i == attempts || !retryableNetErr(err) {
				return lastErr
			}
			time.Sleep(cli.Retry.backoff(i))
			continue
		}

		func() {
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
				lastErr = &HTTPError{Code: resp.StatusCode, Body: string(b)}
				return
			}
			if v == nil {
				io.Copy(io.Discard, resp.Body)
				lastErr = nil
				return
			}
			if de := json.NewDecoder(resp.Body).Decode(v); de != nil {
				lastErr = fmt.Errorf("decode json: %w", de)
				return
			}
			lastErr = nil
		}()

		if lastErr == nil {
			if cli.Log != nil {
				cli.Log.Info("http.ok", "code", resp.StatusCode)
			}
			if cli.Br != nil {
				cli.Br.OnSuccess()
			}
			return nil
		}

		// اگر HTTPError بود و قابل retry، دوباره تلاش کن
		if he, ok := lastErr.(*HTTPError); ok && i < attempts && cli.Retry.Codes[he.Code] {
			time.Sleep(cli.Retry.backoff(i))
			continue
		}

		// در غیر اینصورت برگردون
		return lastErr
	}
	return lastErr
}
