// Package command contains all the commands
package command

import (
	"fmt"
	"minolobot/config"

	"gopkg.in/telegram-bot-api.v4"
)

/* command definition:
help - información general
setactivity - habilita o deshabilita las interacciones
anwswerprob - probabilidad de responder con frases aleatorias
status - estado del bot
*/

// HelpCommand prints the main information
func HelpCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	help := "/help - información general\n" +
		"/setactivity - habilita o deshabilita las interacciones\n" +
		"/anwswerprob - probabilidad de responder con frases aleatorias\n" +
		"/status - estado del bot\n" +
		"pspsps"
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, help)
	bot.Send(msg)
}

// ActivityCommand handles the enabled or diabled status
func ActivityCommand(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	// TODO
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "TO DO")
	bot.Send(msg)
}

// AnswerFrec handles the asignation of a new frequence in the answers
func AnswerFrec(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	// TODO
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "TO DO")
	bot.Send(msg)
}

// Status prints the internat status of the bot
func Status(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	statusData := fmt.Sprintf("Answer: %d%%\nInteraction Enabled: %v\n",
		config.PercentAnswer, config.Enabled)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, statusData)
	bot.Send(msg)
}
