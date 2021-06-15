package util

import "strings"

func LowerTitle(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
