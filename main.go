package main

import (
	"log"
	"math/rand"
	"minolobot/command"
	"minolobot/config"
	"minolobot/interaction"
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
	doPole := interaction.StartPoleLogic()
	listTasks := []func(string) string{interaction.CheckPs, interaction.CheckCion,
		interaction.AnswerDeGoma}
	//////////////////
	// starts routine
	//////////////////
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			// command processing
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "help":
					command.HelpCommand(bot, update)
				case "setactivity":
					command.ActivityCommand(bot, update)
				case "anwswerprob":
					command.AnswerFrec(bot, update)
				case "status":
					command.Status(bot, update)
				}
			}
			// if activity is not enabled just tries to receive the commands
			if !config.Enabled {
				continue
			}
			// pattern processing
			s := string(update.Message.Text)
			for _, task := range listTasks {
				if modString := task(s); modString != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, modString)
					bot.Send(msg)
					continue
				}
			}
			// random answer
			if rand.Intn(99) < config.PercentAnswer {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, interaction.Reply())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		// TODO find way to add generic chat id
		case <-doPole:
			// msg := tgbotapi.NewMessage(12345678, "pole")
			// bot.Send(msg)
		default:
			// TODO which id to do pole?
			// actualTime := time.Now()
			// if actualTime.Hour() == 0 && actualTime.Minute() == 0 {
			// 	msg := tgbotapi.NewMessage(ADD ID HERE, "pole")
			// 	go func(){
			// 		wait := rand.Intn(5000)
			// 		<-time.After(wait * time.Millisecond)
			// 		bot.Send(msg)
			// 	}
			// }
		}
	}
}
