package core

import (
    "net/http"
    "time"
)

type DriverConfig struct {
    Timeout         time.Duration
    Retry           RetryPolicy
    BreakerThresh   int
    BreakerCooldown time.Duration
    HTTPClient      *http.Client
    Logger          Logger
}

func collectDriverConfig(opts ...Option) DriverConfig {
    o := options{}
    for _, fn := range opts { fn(&o) }
    return DriverConfig{
        Timeout: o.Timeout,
        Retry: o.RetryPolicy,
        BreakerThresh: o.BreakerThresh,
        BreakerCooldown: o.BreakerCooldown,
        HTTPClient: o.HTTPClient,
        Logger: o.Logger,
    }
}
