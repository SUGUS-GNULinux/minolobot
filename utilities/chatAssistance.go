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

// Package command contains all the commands
package utilities

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

// Check if the user that send the message is admin
func IsChatAdminUser(bot *tgbotapi.BotAPI, m tgbotapi.Message) bool {
	if m.Chat.IsPrivate() {
		return true
	}

	chatConfig := tgbotapi.ChatConfig{
		ChatID: m.Chat.ID,
	}
	members, err := bot.GetChatAdministrators(chatConfig)

	if err != nil {
		log.Println(err)
		return false
	}

	for _, member := range members {
		if member.User.ID == m.From.ID {
			return true
		}
	}
	return false
}
