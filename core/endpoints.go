package core

// MergeDefaults: دیفالت‌های درایور را با overrideهای کاربر ادغام می‌کند.
func MergeDefaults(def, over map[string]string) map[string]string {
    out := map[string]string{}
    for k, v := range def { out[k] = v }
    for k, v := range over {
        if v == "" { continue }
        out[k] = v
    }
    return out
}
