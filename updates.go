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

// Users is the data structure that contains all the users and their informations
var (
	Users = make(map[int64]*structs.User)
)

// Run the core of the bot
func run(utils types.Utils, data types.Data) {
	// Loop over the updates
	for update := range data.Updates {
		// Save the time of the update reading (more precise than the time of the message)
		curTime := time.Now()

		//Log Update
		utils.Logger.WithFields(logrus.Fields{
			"updID": update.UpdateID,
			"updFr": update.FromChat().Title,
			"updBy": update.SentFrom().UserName,
			"updAt": curTime.Format(utils.TimeFormat),
		}).Debug("Update received")

		// Check the type of the update
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{}).Info("CallbackQuery received")
			// TODO: Manage CallbackQuery
		}
		if update.Message != nil {
			// Log Message
			utils.Logger.WithFields(logrus.Fields{
				"usrFrom": update.Message.From.UserName,
				"msgText": update.Message.Text,
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
				"curTime": curTime.Format(utils.TimeFormat),
			}).Info("Message received")

			// TODO: Rework better this timing system
			eventKey := update.Message.Time().Format("15:04")

			// Check if the message is a command (and ignore other actions)
			if update.Message.IsCommand() {
				manageCommands(update, utils, data, curTime, eventKey)
				continue
			}

			// Check if the message is a valid event
			if event, ok := events.Events.Map[eventKey]; ok && string(eventKey) == update.Message.Text {
				// Log Event message
				utils.Logger.WithFields(logrus.Fields{
					"evnt": update.Message.Text,
					"user": update.Message.From.UserName,
				}).Debug("Event validated")

				// Check if the user has already partecipated
				if event.Activation == nil {
					// Add the user to the data structure if they have never participated before
					if _, ok := Users[update.Message.From.ID]; !ok {
						Users[update.Message.From.ID] = structs.NewUser(update.Message.From.ID, update.Message.From.UserName)
					}

					// Check (and eventually update) the user effects
					UpdateUserEffects(update.Message.From.ID)

					// Activate the event and calculate the delay from o' clock
					event.Activate(Users[update.Message.From.ID], curTime, update.Message.Time(), event.Points)
					delay := curTime.Sub(time.Date(event.Activation.ArrivedAt.Year(), event.Activation.ArrivedAt.Month(), event.Activation.ArrivedAt.Day(), event.Activation.ArrivedAt.Hour(), event.Activation.ArrivedAt.Minute(), 0, 0, event.Activation.ArrivedAt.Location()))

					if event.Activation.ArrivedAt.Second() == 59 {
						event.AddEffect(structs.LastChanceBonus)
					}

					// Apply all effects
					effectText := ""
					curEffects := append(event.Effects, Users[update.Message.From.ID].Effects...)
					if len(curEffects) != 0 {
						effectText += " grazie agli effetti:\n"
						for i := 0; i < len(curEffects); i++ {
							if i != len(curEffects)-1 {
								effectText += fmt.Sprintf("%q, ", curEffects[i].Name)
							} else {
								effectText += fmt.Sprintf("%q", curEffects[i].Name)
							}

							switch curEffects[i].Key {
							case "*":
								event.Activation.EarnedPoints *= curEffects[i].Value
							case "+":
								event.Activation.EarnedPoints += curEffects[i].Value
							case "-":
								event.Activation.EarnedPoints -= curEffects[i].Value
							}
						}
					}

					// Respond to the user with event activated informations
					var msg tgbotapi.MessageConfig
					switch {
					case event.Activation.EarnedPoints < -1:
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Activation.EarnedPoints, effectText, delay.Seconds()))
					case event.Activation.EarnedPoints == -1:
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punto per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Activation.EarnedPoints, effectText, delay.Seconds()))
					case event.Activation.EarnedPoints == 0:
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Peccato %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Activation.EarnedPoints, effectText, delay.Seconds()))
					case event.Activation.EarnedPoints == 1:
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punto per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Activation.EarnedPoints, effectText, delay.Seconds()))
					case event.Activation.EarnedPoints > 1:
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Activation.EarnedPoints, effectText, delay.Seconds()))
					}

					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)

					// Log Event activated
					utils.Logger.WithFields(logrus.Fields{
						"actBy": update.Message.From.UserName,
						"actAt": update.Message.Text,
						"dfPts": event.Points,
						"efPts": event.Activation.EarnedPoints,
					}).Debug("Event activated")

					// Add points to the user if they have never participated the event before
					if event.HasPartecipated(update.Message.From.ID) {
						event.Partecipate(Users[update.Message.From.ID], curTime)
						Users[update.Message.From.ID].TotalPoints += event.Activation.EarnedPoints
						Users[update.Message.From.ID].TotalEventPartecipations++
						Users[update.Message.From.ID].TotalEventWins++
					}
				} else {
					// Calculate the delay from o' clock and winner user
					delay := curTime.Sub(time.Date(event.Activation.ArrivedAt.Year(), event.Activation.ArrivedAt.Month(), event.Activation.ArrivedAt.Day(), event.Activation.ArrivedAt.Hour(), event.Activation.ArrivedAt.Minute(), 0, 0, event.Activation.ArrivedAt.Location()))
					delta := curTime.Sub(event.Activation.ActivatedAt)

					// Respond to the user with event already activated informations
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs.", event.Activation.ActivatedBy, delta.Seconds(), delay.Seconds()))
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)

					// Log Event already activated
					utils.Logger.WithFields(logrus.Fields{
						"actBy": event.Activation.ActivatedBy,
						"actAt": event.Activation.ActivatedAt.Format(utils.TimeFormat),
						"delta": delta,
						"delay": delay,
					}).Debug("Event already activated")

					// Add the user to the data structure if they have never participated before
					if _, ok := Users[update.Message.From.ID]; !ok {
						Users[update.Message.From.ID] = structs.NewUser(update.Message.From.ID, update.Message.From.UserName)
					}
					// Add partecipations to the user if they have never participated the event before
					if event.HasPartecipated(update.Message.From.ID) {
						event.Partecipate(Users[update.Message.From.ID], curTime)
						Users[update.Message.From.ID].TotalEventPartecipations++
					}
				}

				// Save the users file with updated Users data structure
				file, err := json.MarshalIndent(Users, "", " ")
				if err != nil {
					utils.Logger.WithFields(logrus.Fields{
						"err":  err,
						"note": "preoccupati",
					}).Error("Error while marshalling data")
					utils.Logger.Error(Users)
				}
				err = os.WriteFile("files/users.json", file, 0644)
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

func UpdateUserEffects(userID int64) {
	// Generate the ranking
	ranking := make([]Rank, 0)
	for _, u := range Users {
		if u != nil {
			ranking = append(ranking, Rank{u.UserName, u.TotalPoints, u.TotalEventPartecipations})
		}
	}

	// Sort the ranking by points (and partecipations if points are equal)
	sort.Slice(
		ranking,
		func(i, j int) bool {
			if ranking[i].Points == ranking[j].Points {
				return ranking[i].Partecipations < ranking[j].Partecipations
			}
			return ranking[i].Points > ranking[j].Points
		},
	)

	// Get interval from ranking leader
	leaderPoints := ranking[0].Points
	userPoints := 0
	for _, rank := range ranking {
		if rank.Username == Users[userID].UserName {
			userPoints = rank.Points
		}
	}
	interval := leaderPoints - userPoints

	//Remove the Comeback effect
	Users[userID].RemoveEffects(structs.ComebackBonus1, structs.ComebackBonus2, structs.ComebackBonus3)
	switch {
	case interval >= 20 && interval < 50:
		//Add the +1 Comeback effect
		Users[userID].AddEffects(structs.ComebackBonus1)
	case interval >= 50 && interval < 80:
		//Add the +2 Comeback effect
		Users[userID].AddEffects(structs.ComebackBonus2)
	case interval >= 80:
		//Add the +3 Comeback effect
		Users[userID].AddEffects(structs.ComebackBonus3)
	}
}
