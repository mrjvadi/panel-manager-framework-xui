package core

import (
	"regexp"
	"strconv"
	"strings"
)

var nonNum = regexp.MustCompile(`[^0-9]+`)

func parseVersion(v string) (nums []int, pre string, hasPre bool) {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, "-", 2)
	base := parts[0]
	if len(parts) > 1 {
		pre = parts[1]
		hasPre = true
	}
	segs := strings.Split(base, ".")
	for _, s := range segs {
		s = nonNum.ReplaceAllString(s, "")
		if s == "" {
			nums = append(nums, 0)
			continue
		}
		n, _ := strconv.Atoi(s)
		nums = append(nums, n)
	}
	for len(nums) < 3 {
		nums = append(nums, 0)
	}
	return
}

func compareVersionStr(a, b string) int {
	an, apre, ahas := parseVersion(a)
	bn, bpre, bhas := parseVersion(b)
	for i := 0; i < len(an) && i < len(bn); i++ {
		if an[i] < bn[i] {
			return -1
		}
		if an[i] > bn[i] {
			return 1
		}
	}
	// اگر اعداد برابر باشند: pre-release از release کوچک‌تر است
	if ahas && !bhas {
		return -1
	}
	if !ahas && bhas {
		return 1
	}
	// هر دو pre-release: مقایسهٔ سادهٔ رشته‌ای کفایت می‌کند
	if ahas && bhas {
		if apre < bpre {
			return -1
		}
		if apre > bpre {
			return 1
		}
	}
	return 0
}
