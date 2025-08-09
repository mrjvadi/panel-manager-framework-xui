# Panel Manager Framework (X-UI / Marzban) — Go

**Module:** `github.com/mrjvadi/panel-manager-framework-xui`

- Manager/Panel context scopes + Panel.XUI shortcut
- Typed X-UI DTOs + `CloneInbound`, `UpdateInbound`
- Retry (exponential backoff) + circuit breaker
- Pluggable logging (`core.Logger`, adapters for `slog` & `log`)
- Binder for map→struct
- Default endpoint discovery for X-UI forks (overridable via `PanelSpec.Endpoints`)

## Install
```bash
go get github.com/mrjvadi/panel-manager-framework-xui@latest
```

## Quickstart
See `examples/quickstart/main.go` and `examples/xui_clone/main.go`.

## Tests
همهٔ تست‌ها زیر پوشهٔ `./tests` قرار دارند:
```bash
go mod tidy
go test ./...
```

## CI
یک GitHub Action ساده برای build/test زیر `.github/workflows/ci.yml` اضافه شده.
