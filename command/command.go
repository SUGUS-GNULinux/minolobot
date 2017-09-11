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

// Package command contains all the commands
package command

import (
	"fmt"
	"strconv"

	"github.com/SUGUS-GNULinux/minolobot/config"

	"gopkg.in/telegram-bot-api.v4"
	"log"
)

/*
command definition:
help - información general
enable - habilita o deshabilita las interacciones
answer - probabilidad de responder con frases aleatorias
status - estado del bot
*/

const NoCommand string = "hola como uso un comando? unsaludogracias 😂"
const InternalError string = "Uy, algún sugusiano me ha programado mal y acabo de dar un pete interno.\n" +
								"Activando protocolo dedo oreja"
const OnlyAdmin string = "acaso eres admin? Que *P* y que *S*" // Need ParseMode Markdown
const HackingAttempt string = "acabas de disipar todas mis dudas, definitivamente eres tonto"

func init() {
	CommandsFilter = make(map[string]CommandFilter)

	// Maping functions with corresponding command
	addCommand("answer", true, AnswerFreq)
	addCommand("enable", true, EnabledCommand)
	addCommand("help", false, HelpCommand)
	addCommand("pole", true, PoleCommand)
	addCommand("status", true, Status)
}

// HelpCommand prints the main information
func HelpCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	help := "Todas las configuraciones persistiran *únicamente* en el *chat actual*\n" +
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
	msg.ParseMode = "MARKDOWN"
	bot.Send(msg)
}

// EnabledCommand handles the enabled or disabled status
func EnabledCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	var err error
	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")

	switch u.Message.CommandArguments() {
	case "true":
		msg.Text = "ey b0ss"
		err = config.EnabledChatConfig(true, nil, chatID)
	case "false":
		msg.Text = "ye"
		err = config.EnabledChatConfig(false, nil, chatID)
	default:
		msg.Text = NoCommand
	}
	if err != nil {
		log.Println(err)
		msg.Text = InternalError
	}
	bot.Send(msg)
}

// AnswerFreq handles the asignation of a new frequence in the answers
func AnswerFreq(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	value, err := strconv.Atoi(u.Message.CommandArguments())

	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	if err != nil || err == nil && (value > 100 || value < 0) {
		msg.Text = NoCommand
	} else if chatConfig, err := config.FindChatConfig(chatID); err == nil {
		// Do update
		chatConfig.PercentAnswer = value
		err = config.UpdateChatConfig(chatID, chatConfig)
		if err == nil {
			msg.Text = "omgggggg actualizaçao"
		} else {
			log.Println(err)
			msg.Text = InternalError
		}
	} else {
		log.Println(err)
		msg.Text = InternalError
	}
	bot.Send(msg)
}

// PoleCommand handles the enabled or disabled pole
func PoleCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	chatConfig, err := config.FindChatConfig(chatID)
	if err != nil {
		log.Println(err)
		msg.Text = InternalError
	} else {
		switch u.Message.CommandArguments() {
		case "true":
			msg.Text = "👏"
			chatConfig.Pole = true
		case "false":
			msg.Text = "trozo de mierda, tan cansado estás de mis poles? \n😣"
			chatConfig.Pole = false
		default:
			msg.Text = NoCommand
		}
		err = config.UpdateChatConfig(chatID, chatConfig)
		if err != nil {
			msg.Text = InternalError
		}
	}
	bot.Send(msg)
}

// Status prints the internat status of the bot
func Status(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	chatConfig, _ := config.FindChatConfig(chatID)
	statusData := fmt.Sprintf("Answer: %d%%\nInteraction Enabled: %v\nPole Enabled: %v\n",
		chatConfig.PercentAnswer,
		chatConfig.Enabled,
		chatConfig.Pole)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, statusData)
	bot.Send(msg)
}
