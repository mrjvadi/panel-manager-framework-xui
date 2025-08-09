package core

import "time"

// Exported (قبلاً chooseTimeout بود)
func ChooseTimeout(x, def time.Duration) time.Duration {
	if x > 0 {
		return x
	}
	return def
}
