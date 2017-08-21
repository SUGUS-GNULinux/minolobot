package utilities

import (
	"regexp"
	"strings"
)

var regexpSpaces *regexp.Regexp

func init() {
	space := "([ ]+)"
	regexpSpaces = regexp.MustCompile(space)
}

// AnyInSliceIntoString check if the string contains any string from the slice
func AnyInSliceIntoString(a string, list []string) (string, bool) {
	for _, b := range list {
		if strings.Contains(a, b) {
			return b, true
		}
	}
	return "", false
}

// DeleteRedundantSpaces if the string contain more than a space, it replace it by 1 space
func DeleteRedundantSpaces(a string) string {
	res := regexpSpaces.ReplaceAllString(a, " ")
	return res
}
