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

package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/SUGUS-GNULinux/minolobot/command"
	"github.com/SUGUS-GNULinux/minolobot/config"
	"github.com/SUGUS-GNULinux/minolobot/interaction"

	"gopkg.in/telegram-bot-api.v4"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	} else {
		config.BotName = "@" + bot.Self.UserName
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	// needs time to get all the data
	time.Sleep(time.Millisecond * 500)
	// flush the updates when the bot starts
	for len(updates) != 0 {
		<-updates
	}
	// gets the channel to receive the pole signal
	listTasks := []func(string) string{interaction.CheckPs, interaction.CheckCion,
		interaction.AnswerDeGoma}
	//////////////////
	// starts routines
	//////////////////
	go func(bot *tgbotapi.BotAPI) {
		doPole := interaction.StartPoleLogic()
		// every pole signal sends pole message to every registered group
		for range doPole {
			chatsConfig, _ := config.FindAllChatConfigWithId()
			for id, chatConfig := range chatsConfig {
				if chatConfig.Enabled && chatConfig.Pole {
					msg := tgbotapi.NewMessage(id, "pole")
					bot.Send(msg)
				}
			}
		}
	}(bot)

nextUpdate:
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID

		// chat id registration and configuration
		chatConfig, registered := config.FindChatConfig(chatID)
		if registered != nil {
			chatConfig = config.CreateChatConfig(chatID, update.Message.Chat.IsGroup())
			interaction.WelcomeInGroup(bot, update)
		}
		// command processing
		if update.Message.IsCommand() {
			ok := command.AnalyzeCommand(bot, update)
			if ok {
				continue nextUpdate
			}
		}

		// who in sugus
		interaction.Who(bot, update)
		// detection of ignore not enable conditions
		mentionOrPrivate := strings.Contains(update.Message.Text, config.BotName) ||
			update.Message.Chat.IsPrivate()
		// if activity is not enabled just tries to receive the commands.
		// if the message starts with @<botName> or a private msg,
		// it ignores the disabled state.
		if !chatConfig.Enabled && !mentionOrPrivate {
			continue
		}

		if interaction.CheckDisable(bot, update) {
			continue nextUpdate
		}

		// pattern processing
		s := string(update.Message.Text)
		for _, task := range listTasks {
			// if we get content to send from an interaction generator function we send it
			if modString := task(s); modString != "" {
				msg := tgbotapi.NewMessage(chatID, modString)
				msg.ReplyToMessageID = update.Message.MessageID
				msg.ParseMode = "MARKDOWN"
				bot.Send(msg)
				continue nextUpdate
			}
		}
		// random answer
		if rand.Intn(99) < chatConfig.PercentAnswer {
			msg := tgbotapi.NewMessage(chatID, interaction.Reply())
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
