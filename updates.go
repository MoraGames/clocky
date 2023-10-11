package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func run(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
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
