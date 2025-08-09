package core

import "time"

func chooseTimeout(x, def time.Duration) time.Duration {
    if x > 0 { return x }
    return def
}
