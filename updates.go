package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Rank struct {
	Username string
	Points   int
}

var (
	Users = make(map[int64]*structs.User)
)

func run(utils types.Utils, data types.Data) {
	file, err := os.ReadFile("users.json")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err":  err,
			"note": "preoccupati moderatamente",
		}).Error("Error while reading data")
	}

	err = json.Unmarshal([]byte(file), &Users)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err":  err,
			"note": "preoccupati ma poco",
		}).Error("Error while unmarshalling data")
	}

	for update := range data.Updates {
		curTime := time.Now()
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{}).Info("CallbackQuery received")
			//TODO: Manage CallbackQuery
		}
		if update.Message != nil {
			//TODO: Rework better this timing system
			eventKey := events.NewEventKey(update.Message.Time())

			utils.Logger.WithFields(logrus.Fields{
				"usrFrom": update.Message.From.UserName,
				"msgText": update.Message.Text,
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
				"curTime": curTime.Format(utils.TimeFormat),
			}).Info("Message received")

			if update.Message.IsCommand() && update.Message.Command() == "ping" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"msg": update.Message.Text,
				}).Debug("Ping-Pong sent")
			}

			if data.Status == "online" {
				if update.Message.IsCommand() {
					manageCommands(update, utils, data, curTime, eventKey)
				} else if event, ok := events.Events[eventKey]; ok && string(eventKey) == update.Message.Text {
					utils.Logger.WithFields(logrus.Fields{
						"evnt": update.Message.Text,
						"user": update.Message.From.UserName,
					}).Debug("Event validated")
					if !event.Activated {
						event.Activated = true
						event.ActivatedBy = update.Message.From.UserName
						event.ActivatedAt = curTime
						event.ArrivedAt = update.Message.Time()
						delay := event.ActivatedAt.Sub(update.Message.Time())
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! +%v punti per te.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, delay.Seconds()))
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)
						// LOG IMPORTANTI PER VANO -----------------------------------------------
						utils.Logger.WithFields(logrus.Fields{
							"actBy": update.Message.From.UserName,
							"actAt": update.Message.Text,
							"point": event.Points,
							"delay": delay,
						}).Debug("Event activated")

						//give points to the user
						if _, ok := Users[update.Message.From.ID]; !ok {
							Users[update.Message.From.ID] = structs.NewUser(update.Message.From.UserName)
						}
						if !event.Partecipations[update.Message.From.ID] {
							Users[update.Message.From.ID].TotalPoints += event.Points
							Users[update.Message.From.ID].TotalEventPartecipations++
							Users[update.Message.From.ID].TotalEventWins++
						}
					} else {
						delay := curTime.Sub(event.ArrivedAt)
						delta := curTime.Sub(event.ActivatedAt)
						utils.Logger.WithFields(logrus.Fields{
							"evTStr": string(eventKey) + ":00",
							"evTime": event.ArrivedAt,
							"dalay":  delay,
							"dlySec": delay.Seconds(),
							"dlyMil": delay.Milliseconds(),
							"dlyMic": delay.Microseconds(),
							"dlyNan": delay.Nanoseconds(),
						}).Trace("Dalay calculated")
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs", event.ActivatedBy, delta.Seconds(), delay.Seconds()))
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)
						// LOG IMPORTANTI PER VANO -----------------------------------------------
						utils.Logger.WithFields(logrus.Fields{
							"actBy": event.ActivatedBy,
							"actAt": event.ActivatedAt,
							"delta": delta,
							"delay": delay,
						}).Debug("Event already activated")
						if _, ok := Users[update.Message.From.ID]; !ok {
							Users[update.Message.From.ID] = structs.NewUser(update.Message.From.UserName)
						}
						if !event.Partecipations[update.Message.From.ID] {
							Users[update.Message.From.ID].TotalEventPartecipations++
						}
					}

					//set that the user has already partecipated
					events.Events[eventKey].Partecipations[update.Message.From.ID] = true

					file, err := json.MarshalIndent(Users, "", " ")
					if err != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err":  err,
							"note": "preoccupati",
						}).Error("Error while marshalling data")
						utils.Logger.Error(Users)
					}
					err = os.WriteFile("./users.json", file, 0644)
					if err != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err":  err,
							"note": "preoccupati tanto",
						}).Error("Error while writing data")
						utils.Logger.Error(Users)
					}
				}
			}
		}
	}
}
