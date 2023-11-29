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

// Users is the data structure that contains all the users and their informations
var (
	Users = make(map[int64]*structs.User)
)

// Run the core of the bot
func run(utils types.Utils, data types.Data) {
	// Read the users file and load in Users data structure
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

	// Loop over the updates
	for update := range data.Updates {
		// Save the time of the update reading (more precise than the time of the message)
		curTime := time.Now()

		// Check the type of the update
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{}).Info("CallbackQuery received")
			// TODO: Manage CallbackQuery
		}
		if update.Message != nil {
			// TODO: Rework better this timing system
			eventKey := events.NewEventKey(update.Message.Time())

			// Allow only commands if the bot status is offline
			if update.Message.IsCommand() {
				manageCommands(update, utils, data, curTime, eventKey)
			}

			if data.Status == "online" || isAdmin(update.Message.From, utils) {
				// Log Message
				utils.Logger.WithFields(logrus.Fields{
					"usrFrom": update.Message.From.UserName,
					"msgText": update.Message.Text,
					"msgTime": update.Message.Time().Format(utils.TimeFormat),
					"curTime": curTime.Format(utils.TimeFormat),
				}).Info("Message received")

				// Check if the message is a valid event
				if event, ok := events.Events[eventKey]; ok && string(eventKey) == update.Message.Text {
					// Log Event message
					utils.Logger.WithFields(logrus.Fields{
						"evnt": update.Message.Text,
						"user": update.Message.From.UserName,
					}).Debug("Event validated")

					// Check if the user has already partecipated
					if !event.Activated {
						// Activate the event amd calculate the delay from o' clock
						event.Activated = true
						event.ActivatedBy = update.Message.From.UserName
						event.ActivatedAt = curTime
						event.ArrivedAt = update.Message.Time()
						delay := event.ActivatedAt.Sub(update.Message.Time())
						eventTime, err := time.Parse(utils.TimeFormat, string(eventKey))
						if err != nil {
							utils.Logger.WithFields(logrus.Fields{
								"err": err,
							}).Error("Error while parsing event time")
						}
						delay2 := event.ActivatedAt.Sub(eventTime)

						// Respond to the user with event activated informations
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! +%v punti per te.\nHai impiegato +%vs (o forse +%v)", update.Message.From.UserName, event.Points, delay.Seconds(), delay2.Seconds()))
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)

						// Log Event activated
						utils.Logger.WithFields(logrus.Fields{
							"actBy": update.Message.From.UserName,
							"actAt": update.Message.Text,
							"point": event.Points,
						}).Debug("Event activated")

						// Add the user to the data structure if they have never participated before
						if _, ok := Users[update.Message.From.ID]; !ok {
							Users[update.Message.From.ID] = structs.NewUser(update.Message.From.UserName)
						}

						// Add points to the user if they have never participated the event before
						if !event.Partecipations[update.Message.From.ID] {
							Users[update.Message.From.ID].TotalPoints += event.Points
							Users[update.Message.From.ID].TotalEventPartecipations++
							Users[update.Message.From.ID].TotalEventWins++
						}
					} else {
						// Calculate the delay from o' clock and winner user
						delay := curTime.Sub(event.ArrivedAt)
						eventTime, err := time.Parse(utils.TimeFormat, string(eventKey))
						if err != nil {
							utils.Logger.WithFields(logrus.Fields{
								"err": err,
							}).Error("Error while parsing event time")
						}
						delay2 := curTime.Sub(eventTime)
						delta := curTime.Sub(event.ActivatedAt)

						// Respond to the user with event already activated informations
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa (o forse %v).\nHai impiegato +%vs", event.ActivatedBy, delta.Seconds(), delay.Seconds()))
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)

						// Log Event already activated
						utils.Logger.WithFields(logrus.Fields{
							"actBy": event.ActivatedBy,
							"actAt": event.ActivatedAt.Format(utils.TimeFormat),
							"delta": delta,
							"delay": delay,
						}).Debug("Event already activated")

						// Add the user to the data structure if they have never participated before
						if _, ok := Users[update.Message.From.ID]; !ok {
							Users[update.Message.From.ID] = structs.NewUser(update.Message.From.UserName)
						}
						// Add partecipations to the user if they have never participated the event before
						if !event.Partecipations[update.Message.From.ID] {
							Users[update.Message.From.ID].TotalEventPartecipations++
						}
					}

					// Set that the user has already partecipated
					events.Events[eventKey].Partecipations[update.Message.From.ID] = true

					// Save the users file with updated Users data structure
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
