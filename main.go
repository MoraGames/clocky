package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
func main() {
	now := EventKey(time.Now())
	fmt.Println(now.Sintax())
}
*/

func main() {
	bot, err := tgbotapi.NewBotAPI("Please Senpai, put your token in this spot UwU")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			eventKey := toEventKey(update.Message.Time())
			msgTime := time.Now()
			if event, ok := Events[eventKey]; ok && string(eventKey) == update.Message.Text {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				if !event.Activated {
					event.Activated = true
					event.ActivatedBy = update.Message.From.UserName
					event.ActivatedAt = msgTime
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("+%v punti per %v!", event.Points, update.Message.From.UserName))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				} else {
					delta := msgTime.Sub(event.ActivatedAt)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v.\nSei stato più lento di: +%v.%vs", event.ActivatedBy, delta.Seconds(), delta.Milliseconds()))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}
	}
}
