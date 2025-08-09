package core

import "strings"

// compareVersionStr: مقایسه‌ی a و b (رشته‌ای). خروجی: -1 اگر a<b، 0 اگر برابر، 1 اگر a>b
func compareVersionStr(a, b string) int { return compareVersionParts(splitVer(a), splitVer(b)) }

func splitVer(s string) []int {
    s = strings.TrimSpace(s)
    s = strings.TrimPrefix(s, "v")
    s = strings.TrimPrefix(s, "V")
    parts := strings.Split(s, ".")
    out := make([]int, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" { out = append(out, 0); continue }
        n := 0
        for i := 0; i < len(p); i++ {
            if p[i] < '0' || p[i] > '9' { // توقف در اولین غیررقم، مثل "1-rc1"
                break
            }
            n = n*10 + int(p[i]-'0')
        }
        out = append(out, n)
    }
    return out
}

func compareVersionParts(a, b []int) int {
    n := len(a); if len(b) > n { n = len(b) }
    for i := 0; i < n; i++ {
        ai, bi := 0, 0
        if i < len(a) { ai = a[i] }
        if i < len(b) { bi = b[i] }
        if ai < bi { return -1 }
        if ai > bi { return 1 }
    }
    return 0
}
