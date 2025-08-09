package core

import (
    "context"
    "net/http"
    "time"
)

// Hookها (اختیاری)

type Hooks struct {
    OnRequest  func(ctx context.Context, method, url string, req *http.Request)
    OnResponse func(ctx context.Context, code int, dur time.Duration)
}
