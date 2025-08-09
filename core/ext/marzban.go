package ext

// Marzban: اَبَر-اینترفیس مخصوص مرزبان (ترکیب چند افزونه)
// هر درایور مرزبان با implement کردن Subscription + Usage + System به‌صورت خودکار این را هم پیاده می‌کند.

type Marzban interface {
    Subscription
    Usage
    System
}
