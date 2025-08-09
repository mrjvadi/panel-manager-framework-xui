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
    Codes       map[int]bool
}

func DefaultRetryPolicy() RetryPolicy {
    return RetryPolicy{ MaxAttempts: 3, BaseDelay: 150*time.Millisecond, MaxDelay: 1500*time.Millisecond,
        Codes: map[int]bool{408:true,429:true,500:true,502:true,503:true,504:true},
    }
}

func (p RetryPolicy) backoff(n int) time.Duration {
    if n <= 0 { n = 1 }
    d := p.BaseDelay * (1 << (n-1)); if d > p.MaxDelay { d = p.MaxDelay }
    j := time.Duration(rand.Int63n(int64(d)/5))
    if rand.Intn(2)==0 { d -= j } else { d += j }
    return d
}
func retryableNetErr(err error) bool {
    if err == nil { return false }
    if _, ok := err.(net.Error); ok { return true }
    return false
}
