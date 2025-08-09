//go:build integration
// +build integration

package integration

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func getenvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	if n, err := strconv.Atoi(v); err == nil {
		return time.Duration(n) * time.Second
	}
	return def
}

func newHTTPClient(insecure bool, timeout time.Duration) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

func newManagerForHTTP(c *http.Client, reqTimeout time.Duration) *core.Manager {
	return core.New(
		core.WithHTTPClient(c),
		core.WithRequestTimeout(reqTimeout),
	)
}

func ctxTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		d = 10 * time.Second
	}
	return context.WithTimeout(context.Background(), d)
}
