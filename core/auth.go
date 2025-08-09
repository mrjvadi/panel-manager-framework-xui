package core

import (
    "net/http"
    "strings"
)

type Auth interface{ apply(*http.Request) }

type NoAuth struct{}
func (NoAuth) apply(*http.Request) {}

type BasicAuth struct{ Username, Password string }
func (a BasicAuth) apply(r *http.Request) { r.SetBasicAuth(a.Username, a.Password) }

type HeaderToken struct{ Header, Format, Token string }
func (a HeaderToken) apply(r *http.Request) {
    v := a.Token
    if a.Format != "" { v = strings.ReplaceAll(a.Format, "{{token}}", a.Token) }
    r.Header.Set(a.Header, v)
}

type TLS struct{ Insecure bool }
