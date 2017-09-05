package utilities

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

var (
	regexpSpaces        *regexp.Regexp
	ToStringTransformer transform.Transformer
)

func init() {
	space := "([ ]+)"
	regexpSpaces = regexp.MustCompile(space)
}

func init() {

	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}
	ToStringTransformer = transform.Chain(norm.NFD, runes.Remove(runes.Predicate(isMn)), norm.NFC)
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
