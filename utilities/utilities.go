// Copyright 2017-2018 SUGUS GNU/Linux <sugus@us.es>
//
// This file is part of Minolobot.
//
//     Minolobot is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     Minolobot is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with Minolobot.  If not, see <http://www.gnu.org/licenses/>.

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
