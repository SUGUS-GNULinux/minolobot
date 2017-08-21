package utilities

import (
	"fmt"
	"regexp"
	"strings"
)

var regexpSpaces *regexp.Regexp

func init() {
	space := "([ ]+)"
	regexpSpaces = regexp.MustCompile(space)
}

func AnyInSliceIntoString(a string, list []string) (string, bool) {
	for _, b := range list {
		if strings.Contains(a, b) {
			return b, true
		}
	}
	return "", false
}

func DeleteRebundantSpaces(a string) string {
	res := regexpSpaces.ReplaceAllString(a, " ")
	fmt.Println(res)
	return res
}
