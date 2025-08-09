# Panel Manager Framework (X-UI / Marzban) — Go

**Module:** `github.com/mrjvadi/panel-manager-framework-xui`

- Manager/Panel context scopes
- Typed X-UI DTOs + `CloneInbound`
- Retry (exponential backoff) + circuit breaker
- Pluggable logging (`core.Logger`, adapters for `slog` & `log`)
- Binder for map→struct

## Install
```bash
go get github.com/mrjvadi/panel-manager-framework-xui@latest
```

## Quickstart
See `examples/quickstart/main.go`.

## Tests
همهٔ تست‌ها زیر پوشهٔ `./tests` قرار گرفتند:
```bash
go mod tidy
go test ./...
```
تست‌ها از `httptest.Server` استفاده می‌کنند و نیازی به اینترنت ندارند.
