package core

import (
    "crypto/tls"
    "net/http"
    "time"
)

type HTTP struct {
    BaseURL string
    Client  *http.Client
    Retry   RetryPolicy
    Br      *Breaker
    Log     Logger
}

func NewHTTP(base string, insecure bool, timeout time.Duration, c *http.Client) *HTTP {
    if c == nil {
        tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure} }
        c = &http.Client{ Transport: tr, Timeout: timeout }
    } else { c.Timeout = timeout }
    return &HTTP{ BaseURL: base, Client: c, Retry: DefaultRetryPolicy(), Br: NewBreaker(5, 5*time.Second), Log: NoopLogger() }
}
