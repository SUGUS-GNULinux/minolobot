package main

import (
	"log"
	"math/rand"
	"minolobot/command"
	"minolobot/config"
	"minolobot/interaction"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	} else {
		config.BotName = "@" + bot.Self.UserName
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
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
		// TODO find way to add generic chat id
		for _ = range doPole {
			for id := range config.IDList {
				msg := tgbotapi.NewMessage(id, "pole")
				bot.Send(msg)
			}
		}
	}(bot)

nextUpdate:
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// chat id registration
		if update.Message.Chat.IsGroup() && !config.IDList[update.Message.Chat.ID] {
			config.IDList[update.Message.Chat.ID] = true
		}
		// command processing
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				command.HelpCommand(bot, update)
				continue
			case "enable":
				command.ActivityCommand(bot, update)
				continue
			case "anwswer":
				command.AnswerFreq(bot, update)
				continue
			case "status":
				command.Status(bot, update)
				continue
			}
		}
		// detection of ignore not enable conditions
		mentionOrPrivate := strings.Contains(update.Message.Text, config.BotName) ||
			update.Message.Chat.IsPrivate()
		print(mentionOrPrivate)
		// if activity is not enabled just tries to receive the commands.
		// if the message starts with @<botName> or a private msg,
		// it ignores the disabled state.
		if !config.Enabled && !mentionOrPrivate {
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
		if rand.Intn(99) < config.PercentAnswer {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, interaction.Reply())
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
