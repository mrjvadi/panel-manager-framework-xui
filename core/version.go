package core

import (
    "regexp"
    "strconv"
    "strings"
)

var nonNum = regexp.MustCompile(`[^0-9]+`)

func norm(v string) []int {
    v = strings.TrimPrefix(v, "v")
    v = strings.SplitN(v, "-", 2)[0]
    parts := strings.Split(v, ".")
    out := make([]int, 0, len(parts))
    for _, p := range parts {
        p = nonNum.ReplaceAllString(p, "")
        if p == "" { out = append(out, 0); continue }
        n, _ := strconv.Atoi(p)
        out = append(out, n)
    }
    for len(out) < 3 { out = append(out, 0) }
    return out
}

func compareVersionStr(a, b string) int {
    aa := norm(a); bb := norm(b)
    for i := 0; i < len(aa) && i < len(bb); i++ {
        if aa[i] < bb[i] { return -1 }
        if aa[i] > bb[i] { return 1 }
    }
    return 0
}
