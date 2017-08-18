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

package interaction

import (
	"github.com/SUGUS-GNULinux/minolobot/config"
	"gopkg.in/telegram-bot-api.v4"
)

// Send a message when the bot has been added to a group
func WelcomeInGroup(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	if !u.Message.Chat.IsPrivate() {
		text := "Hola, soy " + config.BotName + "\nVamos a *zozobrarnos* juntos !\n" +
			"_(Pssst, si necesitas ayuda usa el comando _/help" + config.BotName + "_)_"

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)
	}
}
