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

						// Add the user to the data structure if they have never participated before
						if _, ok := Users[update.Message.From.ID]; !ok {
							Users[update.Message.From.ID] = structs.NewUser(update.Message.From.UserName)
						}

						// Check (and eventually update) the user effects
						UpdateUserEffects(update.Message.From.ID)

						// Get Event effects
						effectText := ""
						curEffects := append(event.Effects, Users[update.Message.From.ID].Effects...)
						if len(curEffects) != 0 {
							effectText += " grazie agli effetti: "
							for i := 0; i < len(curEffects); i++ {
								if i != len(curEffects)-1 {
									effectText += fmt.Sprintf("%q, ", curEffects[i].Name)
								} else {
									effectText += fmt.Sprintf("%q", curEffects[i].Name)
								}

								switch curEffects[i].Key {
								case "x":
									event.Points *= curEffects[i].Value
								case "+":
									event.Points += curEffects[i].Value
								case "-":
									event.Points -= curEffects[i].Value
								}
							}
						}

						// Respond to the user with event activated informations
						var msg tgbotapi.MessageConfig
						switch {
						case event.Points < -1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, effectText, delay.Seconds()))
						case event.Points == -1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punto per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, effectText, delay.Seconds()))
						case event.Points == 0:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Peccato %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, effectText, delay.Seconds()))
						case event.Points == 1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punto per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, effectText, delay.Seconds()))
						case event.Points > 1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punti per te%v.\nHai impiegato +%vs", update.Message.From.UserName, event.Points, effectText, delay.Seconds()))
						}

						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)

						// Log Event activated
						utils.Logger.WithFields(logrus.Fields{
							"actBy": update.Message.From.UserName,
							"actAt": update.Message.Text,
							"point": event.Points,
						}).Debug("Event activated")

						// Add points to the user if they have never participated the event before
						if !event.Partecipations[update.Message.From.ID] {
							Users[update.Message.From.ID].TotalPoints += event.Points
							Users[update.Message.From.ID].TotalEventPartecipations++
							Users[update.Message.From.ID].TotalEventWins++
						}
					} else {
						// Calculate the delay from o' clock and winner user
						delay := curTime.Sub(time.Date(event.ArrivedAt.Year(), event.ArrivedAt.Month(), event.ArrivedAt.Day(), event.ArrivedAt.Hour(), event.ArrivedAt.Minute(), 0, 0, event.ArrivedAt.Location()))
						delta := curTime.Sub(event.ActivatedAt)

						// Respond to the user with event already activated informations
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs.", event.ActivatedBy, delta.Seconds(), delay.Seconds()))
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

	//Comeback Bonus
	switch {
	case interval < 20:
		//Remove the Comeback effect
		Users[userID].Effects = RemoveUserEffect(Users[userID].Effects, "Comeback Bonus")
	case interval < 50:
		//Add the +1 Comeback effect
		Users[userID].Effects = AddUserEffect(Users[userID].Effects, structs.Effect{Name: "Comeback Bonus", Scope: "User", Key: "+", Value: 1})
	case interval < 80:
		//Add the +2 Comeback effect
		Users[userID].Effects = AddUserEffect(Users[userID].Effects, structs.Effect{Name: "Comeback Bonus", Scope: "User", Key: "+", Value: 2})
	case interval >= 80:
		//Add the +3 Comeback effect
		Users[userID].Effects = AddUserEffect(Users[userID].Effects, structs.Effect{Name: "Comeback Bonus", Scope: "User", Key: "+", Value: 3})
	}
}

func RemoveUserEffect(userEffects []structs.Effect, effectName string) []structs.Effect {
	var newUserEffects []structs.Effect
	for _, e := range userEffects {
		if e.Name != effectName {
			newUserEffects = append(newUserEffects, e)
		}
	}
	return newUserEffects
}

func AddUserEffect(userEffects []structs.Effect, effect structs.Effect) []structs.Effect {
	var newUserEffects []structs.Effect
	for _, e := range userEffects {
		if e.Name != effect.Name {
			newUserEffects = append(newUserEffects, e)
		}
	}
	newUserEffects = append(newUserEffects, effect)
	return newUserEffects
}
