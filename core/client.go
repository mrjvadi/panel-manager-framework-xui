package core

import (
    "crypto/tls"
    "net/http"
    "net/http/cookiejar"
    "time"
)

type HTTP struct {
    BaseURL string
    Client  *http.Client
}

func NewHTTP(base string, insecure bool, timeout time.Duration, c *http.Client) *HTTP {
    if c == nil {
        tr := &http.Transport{ TLSClientConfig: &tls.Config{ InsecureSkipVerify: insecure } }
        c = &http.Client{ Transport: tr, Timeout: timeout }
    }
    jar, _ := cookiejar.New(nil)
    c.Jar = jar
    return &HTTP{ BaseURL: base, Client: c }
}
