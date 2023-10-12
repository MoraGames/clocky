package main

import (
	"fmt"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/MoraGames/types"
)

func run(utils types.Utils, data types.Data) {
	for update := range data.Updates {
		if update.CallbackQuery != nil {
			utils.Log.WithFields(logrus.Fields{
			}).Info("CallbackQuery received")
			//TODO: Manage CallbackQuery
		}
		if update.Message != nil {
			eventKey := events.NewEventKey(update.Message.Time())
			msgTime := time.Now()
			utils.Log.WithFields(logrus.Fields{
				"msgTime": update.Message.Time(),
				"curTime": msgTime,
			}).Info("Message received")
			if event, ok := events.Events[eventKey]; ok && string(eventKey) == update.Message.Text {
				utils.Log.WithFields(logrus.Fields{
					"evnt": update.Message.Text,
					"user": update.Message.From.UserName,
				}).Debug("Event validated")
				if !event.Activated {
					event.Activated = true
					event.ActivatedBy = update.Message.From.UserName
					event.ActivatedAt = msgTime
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("+%v punti per %v!", event.Points, update.Message.From.UserName))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Log.WithFields(logrus.Fields{
						"actBy": update.Message.From.UserName,
						"actAt": update.Message.Text,
						"point": event.Points,
					}).Debug("Event activated")
				} else {
					delta := msgTime.Sub(event.ActivatedAt)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v.\nSei stato più lento di: +%v.%vs", event.ActivatedBy, delta.Seconds(), delta.Milliseconds()))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Log.WithFields(logrus.Fields{
						"actBy": event.ActivatedBy,
						"actAt": event.ActivatedAt,
						"delta": delta,
					}).Debug("Event already activated")
				}
			}
		}
	}
}
