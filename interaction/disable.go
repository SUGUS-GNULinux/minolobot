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
	"github.com/SUGUS-GNULinux/minolobot/command"
	"github.com/SUGUS-GNULinux/minolobot/config"
	"github.com/SUGUS-GNULinux/minolobot/utilities"
	"golang.org/x/text/transform"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
	"time"
)

var (
	timeExpressions = []string{
		"minuto",
		"hora",
	} // Plurals are not necessary
	undefinedTime = map[string]float64{
		"unas": 2,  // time.Hour
		"una":  1,  // time.Hour
		"unos": 18, // time.Minute
		"un":   1,  // time.Minute
	}
	silencExpressions = map[string]time.Duration{
		"no habl":   5 * time.Minute,
		"call":      18 * time.Minute,
		"silenci":   28 * time.Minute,
		"deshabil":  34 * time.Minute,
		"tomar por": 2 * time.Hour,
	}
	silenceExpressionsKeys []string
)

func init() {
	// Create the silenceExpressionsKeys array from the silencExpressions keys
	silenceExpressionsKeys = make([]string, 0)
	for k := range silencExpressions {
		silenceExpressionsKeys = append(silenceExpressionsKeys, k)
	}
}

// Check intention to disable and defer enable
func CheckDisable(bot *tgbotapi.BotAPI, u tgbotapi.Update) (status bool) {
	status = false
	s := u.Message.Text
	s = strings.ToLower(s)
	var d time.Duration
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	msg.ParseMode = "MARKDOWN"

	// Check if the message is for us
	if !((u.Message.ReplyToMessage != nil && u.Message.ReplyToMessage.From.UserName == config.BotName[1:]) ||
		strings.Contains(s, strings.ToLower(config.BotName[1:]))) {
		return
	}

	s, _, err := transform.String(utilities.ToStringTransformer, s)

	if err != nil {
		msg.Text = command.InternalError
		bot.Send(msg)
		return false
	}

	// Check for generic silence expression
	genericSilence, found := utilities.AnyInSliceIntoString(s, silenceExpressionsKeys)
	if !found {
		return
	}

	// Get time and expression from message
	timeExpression, tim, found := getTimeFromMessage(s)
	if found {
		switch timeExpression {
		case "minuto":
			d = time.Minute
		case "hora":
			d = time.Hour
		}
		d = time.Duration(float64(d) * tim)
	} else if timeExpression != "" {
		msg.Text = timeExpression
		bot.Send(msg)
		return true
	} else {
		// Could be a generic silence expression
		d = silencExpressions[genericSilence]
	}

	// Defer bot autoenable
	scheduleFor := time.Now().Local().Add(d)

	// Disable bot
	err = config.EnabledChatConfig(false, &scheduleFor, u.Message.Chat.ID)
	if err != nil {
		msg.Text = command.InternalError
		return true
	}

	msg.Text = "ya me callo *T.T*"
	bot.Send(msg)

	return true
}

func getTimeFromMessage(a string) (string, float64, bool) {
	a = utilities.DeleteRedundantSpaces(a)
	// Iterate among all possible values
	for _, b := range timeExpressions {
		if i := strings.Index(a, b); i >= 0 {
			// If it contains it, we convert the text before the word into a slice
			aSlice := strings.Fields(a[:i])
			timePos := len(aSlice) - 1
			if timePos >= 0 {
				// The last position must contain a numeric value or a undefinedTime
				timeText := strings.Replace(aSlice[len(aSlice)-1], ",", ".", -1)
				timeText = strings.Replace(timeText, "'", ".", -1)
				time, err := strconv.ParseFloat(timeText, 64)
				if err == nil && time > 0 {
					return b, time, true
				} else if err == nil && time < 0 {
					// If time is negative... It's A Trap!! (youtu.be/4F4qzPbcFiA)
					return command.HackingAttempt, 0, false
				} else if time, ok := undefinedTime[aSlice[len(aSlice)-1]]; ok {
					// If not is a numeric value it could be a generic time
					return b, time, true
				}
				// If not, we continue checking the text discarding the previous
			}

			return getTimeFromMessage(a[i+len(b)+1:])
		}
	}
	return "", 0, false
}
