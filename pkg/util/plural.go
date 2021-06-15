package util

import (
	"regexp"
)

var sRegex = regexp.MustCompile(`(^.+)(s|es)$`)

//ToSingular #TODO : special case
func ToSingular(plural string) string {
	return sRegex.ReplaceAllString(plural, "${1}")
}
