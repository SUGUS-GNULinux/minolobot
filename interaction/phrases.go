// Package interaction contains the functions related to interaction with the user
package interaction

import (
	"math/rand"
	"minolobot/config"
	"strings"
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
	if strings.Contains(s, "ps") {
		res = "pspspsps"
	}
	return
}
