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

package command

import (
	"github.com/SUGUS-GNULinux/minolobot/utilities"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

// ChatConfig contains the base configuration for a single user
type CommandFilter struct {
	// Text with which must start the message
	Command string
	// Define is can only be execute by admins
	OnlyAdmin bool
	// Function to execute in positive case
	Func func(bot *tgbotapi.BotAPI, u tgbotapi.Update)
}

var CommandsFilter map[string]CommandFilter

// Analize the command and execute the corresponding function
func AnalyzeCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	c := u.Message.Command()
	actCommand, exist := CommandsFilter[c]

	// Check if command detected
	if !exist {
		return
	}

	// Check only admin
	if u.Message.Chat.IsPrivate() ||
		(!actCommand.OnlyAdmin || utilities.IsChatAdminUser(bot, *u.Message)) {
		actCommand.Func(bot, u)
	} else {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, OnlyAdmin)
		bot.Send(msg)
	}
}

// Add a command to the CommandsFilter
func addCommand(command string, onlyAdmin bool, f func(bot *tgbotapi.BotAPI, u tgbotapi.Update)) {
	_, exist := CommandsFilter[command]
	if exist {
		log.Fatal("Comando existente: ", command)
	}
	CommandsFilter[command] = CommandFilter{
		Command:   command,
		OnlyAdmin: onlyAdmin,
		Func:      f,
	}
}
