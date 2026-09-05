package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Users is the data structure that contains all the users and their informations
var (
	Users        = make(map[int64]*structs.User)
	UserTrackers = make(events.UserTrackersMap)
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

		// Validate the origin of the update
		if updChat.Type != "private" && updChat.ID != App.DefaultChatID {
			utils.Logger.WithFields(fields).Debug("Update ignored due to invalid chat")
			continue
		}

		// Log Update
		utils.Logger.WithFields(fields).Debug("Update received")

		// Check the type of the update
		if update.CallbackQuery != nil {
			utils.Logger.WithFields(logrus.Fields{}).Debug("CallbackQuery received")
			handleButtonComboCallback(update, curTime, utils)
			continue
		}
		if update.Message != nil {
			// Log Message
			utils.Logger.WithFields(logrus.Fields{
				"usrFrom": update.Message.From.UserName,
				"msgText": update.Message.Text,
				"msgTime": update.Message.Time().Format(utils.TimeFormat),
				"curTime": curTime.Format(utils.TimeFormat),
			}).Debug("Message received")

			// Check if the user tracker exists
			if _, exist := UserTrackers[update.Message.From.ID]; !exist {
				UserTrackers[update.Message.From.ID] = events.NewUserTracker(update.Message.From)
			}

			// TODO: Rework better this timing system
			eventKey := update.Message.Time().Format("15:04")

			// Check if the message is a command (and ignore other actions)
			if types.IsCommand(update.Message) {
				manageCommands(update, curTime)
				SaveUserTrackers(utils)
				continue
			}

			// Check if the message is a valid event and if it is enabled
			if event, ok := events.Events.Map[eventKey]; ok && events.IsValidEventMessage(event, update.Message.Text) && event.Enabled {
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
					} else if name := structs.DisplayName(update.Message.From); user.UserName != name {
						// Update the username in case it has changed
						user.UserName = name
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
					if event.ButtonCombo {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Evento speciale!\nPremi il pulsante più velocemente di tutti gli altri.")
						msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Premi qui!", "button_combo:"+event.Name),
							),
						)
					} else {
						switch {
						case event.Activation.EarnedPoints < -1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punti per te%v.\nHai impiegato +%.3fs", user.UserName, event.Activation.EarnedPoints, effectText, delay.Round(time.Millisecond).Seconds()))
						case event.Activation.EarnedPoints == -1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Accidenti %v! %v punto per te%v.\nHai impiegato +%.3fs", user.UserName, event.Activation.EarnedPoints, effectText, delay.Round(time.Millisecond).Seconds()))
						case event.Activation.EarnedPoints == 0:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Peccato %v! %v punti per te%v.\nHai impiegato +%.3fs", user.UserName, event.Activation.EarnedPoints, effectText, delay.Round(time.Millisecond).Seconds()))
						case event.Activation.EarnedPoints == 1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punto per te%v.\nHai impiegato +%.3fs", user.UserName, event.Activation.EarnedPoints, effectText, delay.Round(time.Millisecond).Seconds()))
						case event.Activation.EarnedPoints > 1:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! %v punti per te%v.\nHai impiegato +%.3fs", user.UserName, event.Activation.EarnedPoints, effectText, delay.Round(time.Millisecond).Seconds()))
						}
					}

					msg.ReplyToMessageID = update.Message.MessageID
					sentMessage, err := data.Bot.Send(msg)
					if event.ButtonCombo && err == nil {
						event.SetButtonComboMessageID(sentMessage.MessageID)
						scheduleButtonComboExpiration(event, data.Bot)
					}

					// Log Event activated
					utils.Logger.WithFields(logrus.Fields{
						"actBy": update.Message.From.UserName,
						"actAt": update.Message.Text,
						"dfPts": event.Points,
						"efPts": event.Activation.EarnedPoints,
					}).Debug("Event activated")

					// Add points to the user if they have never participated the event before
					hasPartecipated := event.HasPartecipated(update.Message.From.ID)
					if !hasPartecipated {
						event.Partecipate(user, curTime)
						user.TotalEventPartecipations++
						user.ChampionshipEventPartecipations++
						user.DailyEventPartecipations++
						if !event.ButtonCombo {
							user.TotalPoints += event.Activation.EarnedPoints
							user.TotalEventWins++
							user.ChampionshipPoints += event.Activation.EarnedPoints
							user.ChampionshipEventWins++
							user.DailyPoints += event.Activation.EarnedPoints
							user.DailyEventWins++
						}
					}

					// Update the user in the data structure
					Users[update.Message.From.ID] = user

					// Track: Event message received, validated and participated or won
					UserTrackers[update.Message.From.ID].PushActivity(structs.Activity{
						TelegramTime:          update.Message.Time(),
						ServerReceivingTime:   curTime,
						ServerCompletionTime:  time.Now(),
						Type:                  structs.EventParticipationActivity,
						Message:               update.Message.Text,
						SuccessfulInteraction: !hasPartecipated,
						WinnerUserID:          0,
					})
					if !event.ButtonCombo {
						activity := UserTrackers[update.Message.From.ID].DailyActivities[len(UserTrackers[update.Message.From.ID].DailyActivities)-1].Activities
						activity[len(activity)-1].Type = structs.EventWinActivity
						activity[len(activity)-1].WinnerUserID = update.Message.From.ID
					}
				} else {
					// Calculate the delay from o' clock and winner user
					delay := curTime.Sub(time.Date(event.Activation.ArrivedAt.Year(), event.Activation.ArrivedAt.Month(), event.Activation.ArrivedAt.Day(), event.Activation.ArrivedAt.Hour(), event.Activation.ArrivedAt.Minute(), 0, 0, event.Activation.ArrivedAt.Location()))
					delta := curTime.Sub(event.Activation.ActivatedAt)

					// Respond to the user with event already activated informations
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%.3fs fa.\nHai impiegato +%.3fs.", event.Activation.ActivatedBy.UserName, delta.Round(time.Millisecond).Seconds(), delay.Round(time.Millisecond).Seconds()))
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
					} else if name := structs.DisplayName(update.Message.From); user.UserName != name {
						// Update the username in case it has changed
						user.UserName = name
					}
					if user.TelegramUser == nil {
						AddTelegramUserToExistingUser(update.Message.From)
					}

					// Add partecipations to the user if they have never participated the event before
					hasPartecipated := event.HasPartecipated(update.Message.From.ID)
					if !hasPartecipated {
						event.Partecipate(user, curTime)
						user.TotalEventPartecipations++
						user.ChampionshipEventPartecipations++
						user.DailyEventPartecipations++
					}

					// Update the user in the data structure
					Users[update.Message.From.ID] = user

					// Track: Event message received, validated and participated
					UserTrackers[update.Message.From.ID].PushActivity(structs.Activity{
						TelegramTime:          update.Message.Time(),
						ServerReceivingTime:   curTime,
						ServerCompletionTime:  time.Now(),
						Type:                  structs.EventParticipationActivity,
						Message:               update.Message.Text,
						SuccessfulInteraction: !hasPartecipated,
						WinnerUserID:          event.Activation.ActivatedBy.TelegramID,
					})
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
			} else {
				// Track: Non-event message received (or event disabled)
				UserTrackers[update.Message.From.ID].PushActivity(structs.Activity{
					TelegramTime:          update.Message.Time(),
					ServerReceivingTime:   curTime,
					ServerCompletionTime:  time.Now(),
					Type:                  structs.NonEventActivity,
					Message:               update.Message.Text,
					SuccessfulInteraction: false,
					WinnerUserID:          0,
				})
			}

			SaveUserTrackers(utils)
		}
	}
}

func scheduleButtonComboExpiration(event *events.Event, bot *tgbotapi.BotAPI) {
	delay := time.Until(event.Time.Add(time.Minute))
	if delay < 0 {
		delay = 0
	}
	time.AfterFunc(delay, func() {
		messageID, expired := event.ExpireButtonCombo()
		if !expired || messageID == 0 {
			return
		}
		if _, err := bot.Request(tgbotapi.DeleteMessageConfig{
			ChatID:    App.DefaultChatID,
			MessageID: messageID,
		}); err != nil {
			App.Logger.WithError(err).WithField("event", event.Name).Error("Button combo message not deleted")
		}
	})
}

func handleButtonComboCallback(update tgbotapi.Update, receivedAt time.Time, utils types.Utils) {
	callback := update.CallbackQuery
	if callback == nil || !strings.HasPrefix(callback.Data, "button_combo:") {
		return
	}

	eventName := strings.TrimPrefix(callback.Data, "button_combo:")
	event, found := events.Events.Map[eventName]
	if !found {
		_, _ = App.BotAPI.Request(tgbotapi.NewCallback(callback.ID, "Sfida non disponibile."))
		return
	}

	winnerID, bonus, resolved := event.TryResolveButtonCombo(callback.From.ID)
	if !resolved {
		_, _ = App.BotAPI.Request(tgbotapi.NewCallback(callback.ID, "La sfida è già terminata."))
		return
	}

	winner, found := Users[winnerID]
	if !found {
		_, _ = App.BotAPI.Request(tgbotapi.NewCallback(callback.ID, "Utente non disponibile."))
		return
	}
	winner.TotalPoints += event.Activation.EarnedPoints + bonus
	winner.TotalEventWins++
	winner.ChampionshipPoints += event.Activation.EarnedPoints + bonus
	winner.ChampionshipEventWins++
	winner.DailyPoints += event.Activation.EarnedPoints + bonus
	winner.DailyEventWins++
	Users[winnerID] = winner
	event.Activation.EarnedPoints += bonus

	if _, ok := UserTrackers[winnerID]; !ok {
		UserTrackers[winnerID] = events.NewUserTracker(winner.TelegramUser)
	}
	UserTrackers[winnerID].PushActivity(structs.Activity{
		TelegramTime:          event.Activation.ArrivedAt,
		ServerReceivingTime:   receivedAt,
		ServerCompletionTime:  time.Now(),
		Type:                  structs.EventWinActivity,
		Message:               "button combo",
		SuccessfulInteraction: true,
		WinnerUserID:          winnerID,
	})

	_, _ = App.BotAPI.Request(tgbotapi.NewCallback(callback.ID, "Pulsante registrato!"))
	if callback.Message != nil {
		result := fmt.Sprintf("Complimenti %v! %v punti per te.", winner.UserName, event.Activation.EarnedPoints)
		_, _ = App.BotAPI.Request(tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, result))
	}

	file, err := json.MarshalIndent(Users, "", " ")
	if err == nil {
		if err = os.WriteFile("files/users.json", file, 0644); err != nil {
			utils.Logger.WithError(err).Error("Error while writing Users data")
		}
	}
	SaveUserTrackers(utils)
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
	var leaderPoints int
	if len(ranking) > 0 {
		leaderPoints = ranking[0].Points
	} else {
		// Check if ranking is empty (no users have participated yet)
		// Extra: I'm surprised that this is never happened before. Rankings are empty at the first participation of the championship every time.
		leaderPoints = 0
	}
	userPoints := 0
	for _, rank := range ranking {
		if rank.Username == Users[userID].UserName {
			userPoints = rank.Points
		}
	}
	interval := leaderPoints - userPoints

	//Remove the Comeback effect
	user, ok := Users[userID]
	if !ok {
		// Extra: What? Even there? This crashes was never happened before.
		return
	}
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

func SaveUserTrackers(utils types.Utils) {
	setsFile, err := json.MarshalIndent(UserTrackers, "", "	")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while marshalling Trackers data")
	}
	err = os.WriteFile("files/trackers.json", setsFile, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while writing Trackers data")
	}
}
