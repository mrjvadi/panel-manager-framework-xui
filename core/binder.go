package core

import (
    "bytes"
    "encoding/json"
    "fmt"
)

type Binder struct {
    disallow bool
    src any
}

func Bind() *Binder { return &Binder{} }
func (b *Binder) From(v any) *Binder { b.src = v; return b }
func (b *Binder) DisallowUnknown() *Binder { b.disallow = true; return b }

func (b *Binder) Into(dst any) error {
    if b.src == nil { return fmt.Errorf("binder: nil src") }
    bs, err := json.Marshal(b.src)
    if err != nil { return fmt.Errorf("binder: marshal: %w", err) }
    dec := json.NewDecoder(bytes.NewReader(bs))
    if b.disallow { dec.DisallowUnknownFields() }
    if err := dec.Decode(dst); err != nil { return fmt.Errorf("binder: decode: %w", err) }
    return nil
}

// MapInto: کمک جنریک برای تبدیل مستقیم
func MapInto[T any](src any) (T, error) {
    var dst T
    err := Bind().From(src).Into(&dst)
    return dst, err
}
