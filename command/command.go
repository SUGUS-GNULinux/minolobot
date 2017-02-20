// Package command contains all the commands
package command

import (
	"fmt"
	"minolobot/config"
	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

/*
command definition:
help - información general
enable - habilita o deshabilita las interacciones
answer - probabilidad de responder con frases aleatorias
status - estado del bot
*/

// HelpCommand prints the main information
func HelpCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	help := "/enable - habilita o deshabilita las interacciones añadiendo" +
		" true o false tras el comando.\n" +
		"/answer - probabilidad de responder con frases aleatorias, se define" +
		" añadiendo un numero tras el comando entre 0 y 100\n" +
		"/status - estado del bot\n" +
		"pspsps"
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, help)
	bot.Send(msg)
}

// EnabledCommand handles the enabled or diabled status
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
		msg.Text = "hola como uso un comando? unsaludogracias xd"
	}
	bot.Send(msg)
}

// AnswerFreq handles the asignation of a new frequence in the answers
func AnswerFreq(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	value, err := strconv.Atoi(u.Message.CommandArguments())
	var msg tgbotapi.MessageConfig
	if err != nil || err == nil && (value > 100 || value < 0) {
		msg = tgbotapi.NewMessage(u.Message.Chat.ID,
			"hola como uso un comando? unsaludogracias xd")
	} else {
		chatID := u.Message.Chat.ID
		config.ConfigList[chatID].PercentAnswer = value
		msg = tgbotapi.NewMessage(u.Message.Chat.ID,
			"omgggggg actualizaçao")
	}
	bot.Send(msg)
}

// Status prints the internat status of the bot
func Status(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	chatID := u.Message.Chat.ID
	statusData := fmt.Sprintf("Answer: %d%%\nInteraction Enabled: %v\n",
		config.ConfigList[chatID].PercentAnswer,
		config.ConfigList[chatID].Enabled)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, statusData)
	bot.Send(msg)
}
