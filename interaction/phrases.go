// Copyright 2017 Alejandro Sirgo Rica
// Copyright 2018 Manuel López Ruiz <manuellr.git@gmail.com>
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

// Package interaction contains the functions related to interaction with the user
package interaction

import (
	"math/rand"
	"strings"

	"github.com/SUGUS-GNULinux/minolobot/config"
)

// Reply returns a phrase randomly
func Reply() string {
	return config.Phrases[rand.Intn(len(config.Phrases))]
}

// AnswerDeGoma returns "de goma" as answer to "mojon"
func AnswerDeGoma(s string) (res string) {
	s = strings.ToLower(s)
	if s == "mojon" || s == "mojón" {
		res = "de goma"
	}
	return
}

// CheckCion checks if there is a valid word terminated in "cion" in order to
// modify it and add portuguese termination
func CheckCion(s string) (res string) {
	s = strings.ToLower(s)
	s = strings.Replace(s, "ó", "o", -1)
	if strings.Contains(s, "cion") {
		end := strings.Index(s, "cion") + 4
		if len(s) > end && s[end] != ' ' {
			return
		}
		beg := strings.LastIndex(s[:end], " ") + 1
		cionWord := s[beg:end]
		if ok := config.CionList[cionWord]; ok {
			res = strings.Replace(cionWord, "cion", "çao", 1)
		}
	}
	return
}

// CheckPs returns pspsps to a ps
func CheckPs(s string) (res string) {
	s = strings.ToLower(s)
	s = strings.Replace(s, "https", "", -1)
	if strings.Contains(s, "ps") {
		res = strings.Repeat("ps", 4+rand.Intn(10))
	}
	return
}
