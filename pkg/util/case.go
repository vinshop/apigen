package util

import "strings"

func LowerTitle(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}
