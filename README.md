# Panel Manager Framework (X-UI / Marzban) — Go

**Module:** `github.com/mrjvadi/panel-manager-framework-xui`

فریم‌ورک انعطاف‌پذیر برای مدیریت چند پنل مثل **Marzban** و خانوادهٔ **X-UI** با درایورهای قابل‌افزودن، اجرای موازی، APIهای type-safe، Retry/Breaker داخلی، و اسکোপ کانتکست روی Manager/Panel.

## نصب
```bash
go get github.com/mrjvadi/panel-manager-framework-xui@latest
```

## شروع سریع
```go
mgr := core.New(
    core.WithMaxConcurrency(16),
    core.WithRequestTimeout(10*time.Second),
)

// افزودن پنل‌ها
_ = mgr.AttachByKind(core.PanelSpec{
    ID: "marz-1",
    BaseURL: "https://marzban.example.com",
    Auth: core.BasicAuth{ Username: "admin", Password: "pass" },
    TLS: core.TLS{ Insecure: true },
}, core.DriverMarzban)

_ = mgr.AttachByKind(core.PanelSpec{
    ID: "xui-s-1",
    BaseURL: "https://xui-sanaei.example.com",
    Auth: core.BasicAuth{ Username: "admin", Password: "pass" },
    TLS: core.TLS{ Insecure: true },
}, core.DriverXUISanaei)
```

## اسکوپ کانتکست روی Manager و Panel
```go
// Manager-level scope
r := mgr.Request(core.WithReqTimeoutOpt(8*time.Second))
users, _ := r.UsersAll()

// Panel-level scope
p := mgr.PanelSession("xui-s-1").WithTimeout(5*time.Second)
inbounds, _ := p.Inbounds()
```

## گروه‌بندی و فیلتر
```go
ids := mgr.XUIAll().WhereVersionPrefix("v2.").IDs()
```

## Marzban (Typed)
```go
if mzt, ok := mgr.As[ext.MarzbanTyped]("marz-1"); ok {
    sys, _ := mzt.SystemInfoTyped(ctx)
    fmt.Println("marzban version:", sys.Version)
}
```

## X-UI (Typed + Clone)
```go
if xt, ok := mgr.As[ext.XUITyped]("xui-s-1"); ok {
    // کلون با گزینه‌ها
    inb, _ := xt.CloneInboundTyped(ctx, 123, xdto.CloneInboundOptions{})
}

// شورتکات بدون ctx از روی پنل
p := mgr.PanelSession("xui-s-1")
inb2, _ := p.XUI().CloneInboundWithPort(123, 24443)
```

## Retry و Circuit Breaker
- `core.WithRetryPolicy(core.RetryPolicy{ MaxAttempts: 3, ... })`
- `core.WithBreaker(threshold, cooldown)`

## Examples
- `examples/features/reqctx.go`
- `examples/features/panel_session_usage.go`
- `examples/features/typed_and_health.go`
- `examples/features/xui_typed.go`
- `examples/features/xui_clone.go`
- `examples/features/panel_xui_clone_shortcut.go`

> در مثال‌ها URL/credentialها placeholder هستند؛ جایگزین کنید.

## تست‌ها
- واحدها در `core/*_test.go`, `drivers/*/*_test.go`  
- برای درایورها از `httptest.Server` استفاده شده و نیاز به شبکهٔ واقعی نیست.

## ساخت
```bash
go test ./...
go build ./...
```

## لایسنس
MIT
