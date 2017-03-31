// Copyright 2017 Alejandro Sirgo Rica
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
		for _ = range doPole {
			for id, chatConfig := range config.ConfigList {
				if chatConfig.Enabled && chatConfig.IsGroup {
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
		// chat id registration
		_, registered := config.ConfigList[update.Message.Chat.ID]
		if !registered {
			config.ConfigList[update.Message.Chat.ID] = config.NewChatConfig(
				update.Message.Chat.IsGroup(),
			)
		}
		// command processing
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				command.HelpCommand(bot, update)
				continue
			case "enable":
				command.EnabledCommand(bot, update)
				continue
			case "answer":
				command.AnswerFreq(bot, update)
				continue
			case "status":
				command.Status(bot, update)
				continue
			}
		}
		// who in sugus
		interaction.Who(bot, update)
		// detection of ignore not enable conditions
		mentionOrPrivate := strings.Contains(update.Message.Text, config.BotName) ||
			update.Message.Chat.IsPrivate()
		// actual chat configuration for the update
		chatConfig := config.ConfigList[update.Message.Chat.ID]
		// if activity is not enabled just tries to receive the commands.
		// if the message starts with @<botName> or a private msg,
		// it ignores the disabled state.
		if !chatConfig.Enabled && !mentionOrPrivate {
			continue
		}
		// pattern processing
		s := string(update.Message.Text)
		for _, task := range listTasks {
			// if we get content to send from an interaction generator function we send it
			if modString := task(s); modString != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, modString)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue nextUpdate
			}
		}
		// random answer
		if rand.Intn(99) < chatConfig.PercentAnswer {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, interaction.Reply())
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
