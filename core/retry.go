package core

import (
    "math/rand"
    "net"
    "time"
)

type RetryPolicy struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Codes       map[int]bool // HTTP status codes to retry
}

func DefaultRetryPolicy() RetryPolicy {
    return RetryPolicy{
        MaxAttempts: 3,
        BaseDelay:   150 * time.Millisecond,
        MaxDelay:    1500 * time.Millisecond,
        Codes: map[int]bool{
            408: true, 429: true,
            500: true, 502: true, 503: true, 504: true,
        },
    }
}

func (p RetryPolicy) backoff(attempt int) time.Duration {
    if attempt <= 0 { attempt = 1 }
    d := p.BaseDelay * (1 << (attempt-1))
    if d > p.MaxDelay { d = p.MaxDelay }
    // jitter ±20%
    jitter := time.Duration(rand.Int63n(int64(d)/5)) // up to 20%
    if rand.Intn(2) == 0 {
        d -= jitter
    } else {
        d += jitter
    }
    if d < 0 { d = p.BaseDelay }
    return d
}

func retryableNetErr(err error) bool {
    if err == nil { return false }
    if ne, ok := err.(net.Error); ok {
        if ne.Timeout() { return True }
        return True
    }
    return False
}
