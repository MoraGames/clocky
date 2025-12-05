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

		// Get the update informations
		updID := update.UpdateID
		updAt := curTime.Format(utils.TimeFormat)
		fields := logrus.Fields{
			"updID": updID,
			"updAt": updAt,
		}

		// Get the update.Chat informations (if available)
		updChat := update.FromChat()
		if updChat != nil {
			fields["updChatType"] = updChat.Type
			fields["updChatID"] = updChat.ID
			if updChat.Type == "private" {
				fields["updChatName"] = updChat.UserName
			} else {
				fields["updChatTitl"] = updChat.Title
			}
		}

		// Get the update.User informations (if available)
		updUser := update.SentFrom()
		if updUser != nil {
			fields["updUserID"] = updUser.ID
			if updUser.UserName != "" {
				fields["updUserName"] = updUser.UserName
			}
		}

		//Log Update
		utils.Logger.WithFields(fields).Debug("Update received")

		// Check the type of the update
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{}).Debug("CallbackQuery received")
			// TODO: Manage CallbackQuery
		}
		if update.Message != nil {
			// Log Message
			utils.Logger.WithFields(logrus.Fields{
				"usrFrom": update.Message.From.UserName,
				"msgText": update.Message.Text,
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
				"curTime": curTime.Format(utils.TimeFormat),
			}).Debug("Message received")

			// TODO: Rework better this timing system
			eventKey := update.Message.Time().Format("15:04")

			// Check if the message is a command (and ignore other actions)
			if types.IsCommand(update.Message) {
				manageCommands(update)
				continue
			}

			// Check if the message is a valid event and if it is enabled
			if event, ok := events.Events.Map[eventKey]; ok && string(eventKey) == update.Message.Text && event.Enabled {
				// Log Event message
				utils.Logger.WithFields(logrus.Fields{
					"evnt": update.Message.Text,
					"user": update.Message.From.UserName,
				}).Info("Event validated")

				// Check if the user has already partecipated
				if event.Activation == nil {
					// Add the user to the data structure if they have never participated before
					user, exist := Users[update.Message.From.ID]
					if !exist {
						user = structs.NewUser(update.Message.From)
					} else if user.UserName != update.Message.From.UserName {
						// Update the username in case it has changed
						user.UserName = update.Message.From.UserName
					}
					if user.TelegramUser == nil {
						AddTelegramUserToExistingUser(update.Message.From)
					}

					// Check (and eventually update) the user effects
					UpdateUserEffects(update.Message.From.ID)

					// Activate the event and calculate the delay from o' clock
					event.Activate(user, curTime, update.Message.Time(), event.Points)
					delay := curTime.Sub(time.Date(event.Activation.ArrivedAt.Year(), event.Activation.ArrivedAt.Month(), event.Activation.ArrivedAt.Day(), event.Activation.ArrivedAt.Hour(), event.Activation.ArrivedAt.Minute(), 0, 0, event.Activation.ArrivedAt.Location()))

					if event.Activation.ArrivedAt.Second() == 58 {
						event.AddEffect(structs.LastChanceBonus)
					} else if event.Activation.ArrivedAt.Second() == 59 {
						event.AddEffect(structs.LastChanceBonus2)
					}

					// Apply all effects
					effectText := ""
					curEffects := append(event.Effects, user.Effects...)
					if len(curEffects) != 0 {
						effectText += " grazie agli effetti:\n"
						for i := 0; i < len(curEffects); i++ {
							if i != len(curEffects)-1 {
								effectText += fmt.Sprintf("%q, ", curEffects[i].Name)
							} else {
								effectText += fmt.Sprintf("%q", curEffects[i].Name)
							}

							if (curEffects[i].Name == structs.NoNegative.Name) && event.Activation.EarnedPoints < 0 {
								event.Activation.EarnedPoints = 0
								continue
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
					if !event.HasPartecipated(update.Message.From.ID) {
						event.Partecipate(user, curTime)
						user.TotalPoints += event.Activation.EarnedPoints
						user.TotalEventPartecipations++
						user.TotalEventWins++
						user.ChampionshipPoints += event.Activation.EarnedPoints
						user.ChampionshipEventPartecipations++
						user.ChampionshipEventWins++
						user.DailyPoints += event.Activation.EarnedPoints
						user.DailyEventPartecipations++
						user.DailyEventWins++
					}

					// Update the user in the data structure
					Users[update.Message.From.ID] = user
				} else {
					// Calculate the delay from o' clock and winner user
					delay := curTime.Sub(time.Date(event.Activation.ArrivedAt.Year(), event.Activation.ArrivedAt.Month(), event.Activation.ArrivedAt.Day(), event.Activation.ArrivedAt.Hour(), event.Activation.ArrivedAt.Minute(), 0, 0, event.Activation.ArrivedAt.Location()))
					delta := curTime.Sub(event.Activation.ActivatedAt)

					// Respond to the user with event already activated informations
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs.", event.Activation.ActivatedBy.UserName, delta.Seconds(), delay.Seconds()))
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
					user, exist := Users[update.Message.From.ID]
					if !exist {
						user = structs.NewUser(update.Message.From)
					} else if user.UserName != update.Message.From.UserName {
						// Update the username in case it has changed
						user.UserName = update.Message.From.UserName
					}
					if user.TelegramUser == nil {
						AddTelegramUserToExistingUser(update.Message.From)
					}
					// Add partecipations to the user if they have never participated the event before
					if !event.HasPartecipated(update.Message.From.ID) {
						event.Partecipate(user, curTime)
						user.TotalEventPartecipations++
						user.ChampionshipEventPartecipations++
						user.DailyEventPartecipations++
					}

					// Update the user in the data structure
					Users[update.Message.From.ID] = user
				}

				// Save the users file with updated Users data structure
				file, err := json.MarshalIndent(Users, "", " ")
				if err != nil {
					utils.Logger.WithFields(logrus.Fields{
						"err": err,
						"msg": "Error while marshalling Users data",
					}).Error("Error while marshalling data")
					utils.Logger.Error(Users)
				}
				err = os.WriteFile("files/users.json", file, 0644)
				if err != nil {
					utils.Logger.WithFields(logrus.Fields{
						"err": err,
						"msg": "Error while writing Users data",
					}).Error("Error while writing data")
					utils.Logger.Error(Users)
				}
			}
		}
	}
}

func UpdateUserEffects(userID int64) {
	// Generate the ranking
	ranking := structs.GetRanking(Users, structs.RankScopeChampionship, false)

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
	user := Users[userID]
	user.RemoveEffect(structs.ComebackBonus1)
	user.RemoveEffect(structs.ComebackBonus2)
	user.RemoveEffect(structs.ComebackBonus3)
	user.RemoveEffect(structs.ComebackBonus4)
	user.RemoveEffect(structs.ComebackBonus5)
	switch {
	case interval >= 20 && interval < 40:
		//Add the +1 Comeback effect
		user.AddEffect(structs.ComebackBonus1)
	case interval >= 40 && interval < 60:
		//Add the +2 Comeback effect
		user.AddEffect(structs.ComebackBonus2)
	case interval >= 60 && interval < 80:
		//Add the +3 Comeback effect
		user.AddEffect(structs.ComebackBonus3)
	case interval >= 80 && interval < 100:
		//Add the +4 Comeback effect
		user.AddEffect(structs.ComebackBonus4)
	case interval >= 100:
		//Add the +5 Comeback effect
		user.AddEffect(structs.ComebackBonus5)
	}
	Users[userID] = user
}
