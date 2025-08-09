package core

// Exported wrapper for tests in external 'tests' package.
func TestCompare(a, b string) int { return compareVersionStr(a, b) }
