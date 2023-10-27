package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
				"curTime": curTime.Format(utils.TimeFormat),
			}).Info("Message received")
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "ranking":
					ranking := make([]Rank, 0)
					for _, u := range Users {
						ranking = append(ranking, Rank{u.UserName, u.TotalPoints})
					}
					sort.Slice(ranking, func(i, j int) bool { return ranking[i].Points > ranking[j].Points })

					rankingString := ""
					for i, r := range ranking {
						rankingString += fmt.Sprintf("%v] %v: %v\n", i+1, r.Username, r.Points)
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("La classifica è la seguente:\n\n%v", rankingString))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.Debug("Ranking sent")
				case "stats":
					u := Users[update.Message.From.ID]
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Le tue statistiche sono:\n\nPunti totali: %v\nPartecipazioni totali: %v\nVittorie totali: %v", u.TotalPoints, u.TotalEventPartecipations, u.TotalEventWins))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.Debug("Stats sent")
				}
			} else if event, ok := events.Events[eventKey]; ok && string(eventKey) == update.Message.Text {
				utils.Logger.WithFields(logrus.Fields{
					"evnt": update.Message.Text,
					"user": update.Message.From.UserName,
				}).Debug("Event validated")
				if events.LastEventKey != eventKey {
					events.LastEventKey = eventKey
					event.Activated = true
					event.ActivatedBy = update.Message.From.UserName
					event.ActivatedAt = curTime
					event.ArrivedAt = update.Message.Time()
					retard := event.ActivatedAt.Sub(update.Message.Time())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! +%v punti per te.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, retard.Seconds()))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.WithFields(logrus.Fields{
						"actBy": update.Message.From.UserName,
						"actAt": update.Message.Text,
						"point": event.Points,
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
					retard := curTime.Sub(event.ArrivedAt)
					delta := curTime.Sub(event.ActivatedAt)
					utils.Logger.WithFields(logrus.Fields{
						"evTStr": string(eventKey) + ":00",
						"evTime": event.ArrivedAt,
						"retard": retard,
						"retSec": retard.Seconds(),
						"retMil": retard.Milliseconds(),
						"retMic": retard.Microseconds(),
						"retNan": retard.Nanoseconds(),
					}).Warn("Retard calculated")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs", event.ActivatedBy, delta.Seconds(), retard.Seconds()))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.WithFields(logrus.Fields{
						"actBy": event.ActivatedBy,
						"actAt": event.ActivatedAt,
						"delta": delta,
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
