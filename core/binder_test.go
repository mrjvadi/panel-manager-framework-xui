package core

import "testing"

type dto struct {
    A string `json:"a"`
    B int    `json:"b"`
}

func TestBinderOK(t *testing.T) {
    src := map[string]any{"a":"x","b":3}
    var d dto
    if err := Bind().From(src).Into(&d); err != nil { t.Fatal(err) }
    if d.A != "x" || d.B != 3 { t.Fatalf("bad bind: %#v", d) }
}

func TestBinderDisallowUnknown(t *testing.T) {
    src := map[string]any{"a":"x","b":3,"c":9}
    var d dto
    if err := Bind().From(src).DisallowUnknown().Into(&d); err == nil {
        t.Fatal("expected error for unknown field")
    }
}
