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

package interaction

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var regexpQuestion *regexp.Regexp

const url = `https://sugus.eii.us.es/en_sugus.html`

func init() {
	question := "(?:(?:(?:q|Q)ui(?:e|é)n(?: hay| est(?:a|á))?)|(?:A|a)lguien) " +
		"(?:en|por) (?:s|S)ugus"
	regexpQuestion = regexp.MustCompile(question)
}

// Who identifies when someone
func Who(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	match := regexpQuestion.MatchString(u.Message.Text)
	if match {
		var s string
		res, err := http.Get(url)
		if err == nil {
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				s = extractList(res.Body)
				res.Body.Close()
			} else {
				s = fmt.Sprint("HTTP status code ", res.StatusCode)
			}
		} else {
			s = err.Error()
		}
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, s)
		bot.Send(msg)
	}
}

func extractList(r io.Reader) (res string) {
	whoPage, err := ioutil.ReadAll(r)
	if err == nil {
		begin := bytes.Index(whoPage, []byte("<li>"))
		whoPage = whoPage[begin:]
		end := bytes.Index(whoPage, []byte("</ul>"))
		res = string(whoPage[:end])
		res = strings.Replace(res, "<li>", "", -1)
		res = strings.Replace(res, "</li>", "", -1)
	}
	return
}
