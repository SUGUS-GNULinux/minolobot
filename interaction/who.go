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
