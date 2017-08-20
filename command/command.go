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

// Package command contains all the commands
package command

import (
	"fmt"
	"strconv"

	"github.com/SUGUS-GNULinux/minolobot/config"

	"gopkg.in/telegram-bot-api.v4"
)

/*
command definition:
help - información general
enable - habilita o deshabilita las interacciones
answer - probabilidad de responder con frases aleatorias
status - estado del bot
*/

const NoCommand string = "hola como uso un comando? unsaludogracias 😂"
const OnlyAdmin string = "acaso eres admin? Que *P* y que *S*" // Need ParseMode Markdown

// HelpCommand prints the main information
func HelpCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	help := "Todas las configuraciones persistiran únicamente en el chat actual" +
		"/enable - habilita o deshabilita las interacciones añadiendo" +
		" true o false tras el comando.\n" +
		"/answer - probabilidad de responder con frases aleatorias, se define" +
		" añadiendo un numero tras el comando entre 0 y 100\n" +
		"/pole - habilita o deshabilita la pole añadiendo" +
		" true o false tras el comando.\n" +
		"/status - estado del bot\n" +
		"pspsps"
	dest := int64(u.Message.From.ID) // Avoids answering groups
	msg := tgbotapi.NewMessage(dest, help)
	bot.Send(msg)
}

// EnabledCommand handles the enabled or disabled status
func EnabledCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	switch u.Message.CommandArguments() {
	case "true":
		msg.Text = "ey b0ss"
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].Enabled = true
	case "false":
		msg.Text = "ye"
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].Enabled = false
	default:
		msg.Text = NoCommand
	}
	bot.Send(msg)
}

// AnswerFreq handles the asignation of a new frequence in the answers
func AnswerFreq(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	value, err := strconv.Atoi(u.Message.CommandArguments())
	var msg tgbotapi.MessageConfig
	if err != nil || err == nil && (value > 100 || value < 0) {
		msg = tgbotapi.NewMessage(u.Message.Chat.ID, NoCommand)
	} else {
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].PercentAnswer = value
		msg = tgbotapi.NewMessage(u.Message.Chat.ID,
			"omgggggg actualizaçao")
	}
	bot.Send(msg)
}

// PoleCommand handles the enabled or disabled pole
func PoleCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	switch u.Message.CommandArguments() {
	case "true":
		msg.Text = "👏"
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].Pole = true
	case "false":
		msg.Text = "trozo de mierda, tan cansado estás de mis poles? \n😣"
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].Pole = false
	default:
		msg.Text = NoCommand
	}
	bot.Send(msg)
}

// Status prints the internat status of the bot
func Status(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	statusData := fmt.Sprintf("Answer: %d%%\nInteraction Enabled: %v\nPole Enabled: %v\n",
		config.ConfigList[chatID].PercentAnswer,
		config.ConfigList[chatID].Enabled,
		config.ConfigList[chatID].Pole)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, statusData)
	bot.Send(msg)
}
