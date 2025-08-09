package core

import "testing"

func TestCompareVersionStr(t *testing.T) {
    cases := []struct{ a, b string; want int }{
        {"v1.2.3", "1.2.3", 0},
        {"v1.2.3", "1.2.4", -1},
        {"1.10.0", "1.2.9", 1},
        {"1.2", "1.2.0", 0},
        {"2.0-rc1", "2.0", -1},
    }
    for _, c := range cases {
        got := compareVersionStr(c.a, c.b)
        if got != c.want { t.Fatalf("%s ? %s: got %d want %d", c.a, c.b, got, c.want) }
    }
}
