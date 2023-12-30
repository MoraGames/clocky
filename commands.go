package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Rank is the type used for manage /ranking sorting
type Rank struct {
	Username       string
	Points         int
	Partecipations int
}

// switch for all the commands that the bot can receive
func manageCommands(update tgbotapi.Update, utils types.Utils, data types.Data, curTime time.Time, eventKey string) {
	switch update.Message.Command() {
	case "check":
		// Check actual event infos
		if !isAdmin(update.Message.From, utils) {
			// Respond and log with a message indicating that the user is not authorized to use this command
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
			msg.ReplyToMessageID = update.Message.MessageID
			message, error := data.Bot.Send(msg)
			if error != nil {
				utils.Logger.WithFields(logrus.Fields{
					"err": error,
					"msg": message,
				}).Error("Error while sending message")
			}
			utils.Logger.WithFields(logrus.Fields{
				"usr": update.Message.From.UserName,
				"cmd": update.Message.Command(),
			}).Debug("Unauthorized user")
		} else {
			// Split the command arguments
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")

			if len(cmdArgs) != 1 {
				// Respond with a message indicating that the command arguments are wrong
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /check <events|users|logs>")
				msg.ReplyToMessageID = update.Message.MessageID
				message, error := data.Bot.Send(msg)
				if error != nil {
					utils.Logger.WithFields(logrus.Fields{
						"err": error,
						"msg": message,
					}).Error("Error while sending message")
				}
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"msg": update.Message.Text,
				}).Debug("Wrong command")
			} else {
				// Check if the command argument is events
				switch cmdArgs[0] {
				case "logs":
					// Check the logs data structure
					logTxt, err := os.ReadFile("files/log.txt")
					if err != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": err,
						}).Error("Error while reading files/log.txt")
					}

					// Respond with command executed successfully
					msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "log.txt", Bytes: logTxt})
					msg.Caption = "Log controllati. Ecco lo stato attuale:\n\n"
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /check command sent
					utils.Logger.Debug("Logs checked")
				case "users":
					// Check the logs data structure
					usersJson, err := os.ReadFile("files/users.json")
					if err != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": err,
						}).Error("Error while reading files/users.json")
					}

					// Respond with command executed successfully
					msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "users.json", Bytes: usersJson})
					msg.Caption = "Log controllati. Ecco lo stato attuale:\n\n"
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /check command sent
					utils.Logger.Debug("Users checked")
				case "events":
					// Check the events data structure
					eventsJson, err := json.MarshalIndent(events.Events, "", " ")
					if err != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err":  err,
							"note": "preoccupati",
						}).Error("Error while marshalling Events data")
						utils.Logger.Error(Users)
					}

					// Respond with command executed successfully
					msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "events.json", Bytes: eventsJson})
					msg.Caption = "Eventi controllati. Ecco lo stato attuale:\n\n"
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /check command sent
					utils.Logger.Debug("Events checked")
				default:
					// Respond with a message indicating that the command arguments are wrong
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /check <events|users|logs>")
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}
					utils.Logger.WithFields(logrus.Fields{
						"usr": update.Message.From.UserName,
						"msg": update.Message.Text,
					}).Debug("Wrong command")
				}
			}
		}
	case "credits":
		// Respond with useful information about the project
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The source code, available on GitHub at MoraGames/clockyuwu, is written entirely in GoLang and makes use of the \"telegram-bot-api\" library.\nFor any bug reports or feature proposals, please refer to the GitHub project.\n\nDeveloper:\n- Telegram: @MoraGames\n- Discord: @moragames\n- Instagram: @moragames.dev\n- GitHub: MoraGames\n\nProject:\n- Telegram: @clockyuwu_bot\n- GitHub: MoraGames/clockyuwu\n\nSpecial thanks go to the first testers (as well as players) of the minigame managed by the bot, \"Vano\", \"Ale\" and \"Alex\".")
		msg.ReplyToMessageID = update.Message.MessageID
		message, error := data.Bot.Send(msg)
		if error != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": error,
				"msg": message,
			}).Error("Error while sending message")
		}

		utils.Logger.WithFields(logrus.Fields{
			"message": update.Message.Text,
			"sender":  update.Message.From.UserName,
			"chat":    update.Message.Chat.Title,
		}).Debug("Response to \"/credits\" command sent successfully")
	case "help":
		// Respond with useful information about the working and commands of the bot
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Name: %v\nVersion: %v\n\nThis is a list of all possible commands within the bot:\n\n- /start : Get an introductory message about the bot's features.\n - /help : Get a complete list of all available commands.\n - /ranking : Get the ranking of the current championship.\n - /stats : Get the player's game statistics.\n - /ping : Verify if the bot is running.\n - /credits : Get more informations abount the project.\n\nAdmin's Only:\n - /check : Get more informations about bot status and data.\n - /reset : Force the execution of a specific Reset() function.\n -/update : Update the value of a data structure.", utils.Config.App.Name, utils.Config.App.Version))
		msg.ReplyToMessageID = update.Message.MessageID
		message, error := data.Bot.Send(msg)
		if error != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": error,
				"msg": message,
			}).Error("Error while sending message")
		}

		utils.Logger.WithFields(logrus.Fields{
			"message": update.Message.Text,
			"sender":  update.Message.From.UserName,
			"chat":    update.Message.Chat.Title,
		}).Debug("Response to \"/help\" command sent successfully")
	case "ping":
		// Respond with a "pong" message. Useful for checking if the bot is online
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
		msg.ReplyToMessageID = update.Message.MessageID
		message, error := data.Bot.Send(msg)
		if error != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": error,
				"msg": message,
			}).Error("Error while sending message")
		}

		utils.Logger.WithFields(logrus.Fields{
			"message": update.Message.Text,
			"sender":  update.Message.From.UserName,
			"chat":    update.Message.Chat.Title,
		}).Debug("Response to \"/ping\" command sent successfully")
	case "ranking":
		// Respond with the ranking based on users' points
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

		// Generate the string to send
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ancora nessun utente ha partecipato agli eventi della season.")
		if len(ranking) != 0 {
			leadersPoints := ranking[0].Points
			rankingString := ""
			for i, r := range ranking {
				rankingString += fmt.Sprintf("%v] %v: %v (-%v)\n", i+1, r.Username, r.Points, leadersPoints-r.Points)
			}

			// Send the message
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("La classifica è la seguente:\n\n%v", rankingString))
		}
		msg.ReplyToMessageID = update.Message.MessageID
		message, error := data.Bot.Send(msg)
		if error != nil {
			utils.Logger.WithFields(logrus.Fields{
				"err": error,
				"msg": message,
			}).Error("Error while sending message")
		}

		// Log the /ranking command sent
		utils.Logger.Debug("Ranking sent")
	case "reset":
		// Reset the events or users data structure
		// Check if the user is an bot-admin
		if !isAdmin(update.Message.From, utils) {
			// Respond and log with a message indicating that the user is not authorized to use this command
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
			msg.ReplyToMessageID = update.Message.MessageID
			message, error := data.Bot.Send(msg)
			if error != nil {
				utils.Logger.WithFields(logrus.Fields{
					"err": error,
					"msg": message,
				}).Error("Error while sending message")
			}
			utils.Logger.WithFields(logrus.Fields{
				"usr": update.Message.From.UserName,
				"cmd": update.Message.Command(),
			}).Debug("Unauthorized user")
		} else {
			// Split the command arguments
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")

			// Check if the command arguments are in the form /reset <events|users>
			if len(cmdArgs) != 1 {
				// Respond with a message indicating that the command arguments are wrong
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /reset <events|users>")
				msg.ReplyToMessageID = update.Message.MessageID
				message, error := data.Bot.Send(msg)
				if error != nil {
					utils.Logger.WithFields(logrus.Fields{
						"err": error,
						"msg": message,
					}).Error("Error while sending message")
				}
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"msg": update.Message.Text,
				}).Debug("Wrong command")
			} else {
				// Check if the command argument is events or users
				switch cmdArgs[0] {
				case "events":
					// Reset the events data structure
					events.Events.Reset(
						true,
						&types.WriteMessageData{Bot: data.Bot, ChatID: update.Message.Chat.ID, ReplyMessageID: update.Message.MessageID},
						utils,
					)

					// Respond with command executed successfully
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Eventi resettati")
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /reset command sent
					utils.Logger.Debug("Events resetted")
				case "users":
					// Reset the users data structure
					Users = make(map[int64]*structs.User)

					// Overwrite the files/users.json file with the new (and empty) data structure
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

					// Respond with command executed successfully
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Utenti resettati")
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /reset command sent
					utils.Logger.Debug("Users resetted")
				default:
					// Respond with a message indicating that the command arguments are wrong
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /reset <events|users>")
					msg.ReplyToMessageID = update.Message.MessageID
					message, error := data.Bot.Send(msg)
					if error != nil {
						utils.Logger.WithFields(logrus.Fields{
							"err": error,
							"msg": message,
						}).Error("Error while sending message")
					}

					// Log the /reset command executed in a wrong form
					utils.Logger.WithFields(logrus.Fields{
						"usr": update.Message.From.UserName,
						"cmd": update.Message.Command(),
					}).Debug("Wrong command")
				}
			}
		}
	case "start":
		// Respond with an introduction message for the users of the bot
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			ComposeMessage([]string{
				"%v is a bot that allows you to play a time-wasting game with one or more groups of friends within Telegram groups.",
				"Once the bot is added, the game mainly (but not exclusively) involves sending messages in the \"hh:mm\" format at certain times of the day, in exchange for valuable points.",
				"The person who has earned the most points at the end of the championship will be the new Clocky Champion!\n",
				"Use /help to get a list of all commands or /credits for more information about the project.\n\n- %v, a bot from @MoraGames.",
			}, utils.Config.App.Name, utils.Config.App.Name),
		)
		SendMessage(msg, update, data, utils)
		// Log the command executed successfully
		FinalCommandLog("\"start message sent", update, utils)
		SuccessResponseLog(update, utils)
	case "stats":
		/*
			Description:
				Get the stats of the user who sent the command or of the user specified in the command arguments.

			Forms:
				/stats
				/stats [user]
		*/
		// Check if the command has arguments
		if update.Message.CommandArguments() == "" {
			// Get the user from the Users data structure
			u := Users[update.Message.From.ID]
			// Check (and eventually update) the user effects
			UpdateUserEffects(update.Message.From.ID)
			// Send the message with user's stats
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non hai ancora partecipato a nessun evento.")
			if u != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Le tue statistiche sono:\n\nPunti totali: %v\nPartecipazioni totali: %v\nVittorie totali: %v\nPunti/Partecipazioni: %.2f\nPunti/Vittorie: %.2f\nVittorie/Partecipazioni: %.2f\nVittorie/Sconfitte: %.2f\nEffetti attivi: %v", u.TotalPoints, u.TotalEventPartecipations, u.TotalEventWins, float64(u.TotalPoints)/float64(u.TotalEventPartecipations), float64(u.TotalPoints)/float64(u.TotalEventWins), float64(u.TotalEventWins)/float64(u.TotalEventPartecipations), float64(u.TotalEventWins)/float64(u.TotalEventPartecipations-u.TotalEventWins), u.StringifyEffects()))
			}
			SendMessage(msg, update, data, utils)
			// Log the command executed successfully
			FinalCommandLog("User.Stats sent", update, utils)
			SuccessResponseLog(update, utils)
		} else {
			// Split the command arguments
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")
			if len(cmdArgs) != 1 {
				// Respond with a message indicating that the command arguments are wrong
				cmdSyntax := "/stats [user]"
				SendWrongCommandSyntaxMessage(cmdSyntax, update, data, utils)
				// Log the command failed execution
				FinalCommandLog("Wrong command syntax", update, utils)
			} else {
				// Get and check if the user exists
				username := cmdArgs[0]
				var userKey int64
				var founded bool
				for userID, user := range Users {
					if user.UserName == username {
						founded = true
						userKey = userID
					}
				}
				if !founded {
					// Respond with a message indicating that the user does not exist
					SendEntityNotFoundMessage("Utente", username, update, data, utils)
					// Log the command failed execution
					FinalCommandLog("User not found", update, utils)
				} else {
					// Get the user from the Users data structure
					u := Users[userKey]
					// Check (and eventually update) the user effects
					UpdateUserEffects(userKey)
					// Send the message with user's stats
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v non ha ancora partecipato a nessun evento.", u.UserName))
					if u != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Le statistiche di %v sono:\n\nPunti totali: %v\nPartecipazioni totali: %v\nVittorie totali: %v\nPunti/Partecipazioni: %.2f\nPunti/Vittorie: %.2f\nVittorie/Partecipazioni: %.2f\nVittorie/Sconfitte: %.2f\nEffetti attivi: %v", u.UserName, u.TotalPoints, u.TotalEventPartecipations, u.TotalEventWins, float64(u.TotalPoints)/float64(u.TotalEventPartecipations), float64(u.TotalPoints)/float64(u.TotalEventWins), float64(u.TotalEventWins)/float64(u.TotalEventPartecipations), float64(u.TotalEventWins)/float64(u.TotalEventPartecipations-u.TotalEventWins), u.StringifyEffects()))
					}
					SendMessage(msg, update, data, utils)
					// Log the command executed successfully
					FinalCommandLog("User.Stats sent", update, utils)
					SuccessResponseLog(update, utils)
				}
			}
		}
	case "update":
		/*
			Description:
				Update an event (points, enabled, effects) or user (points, partecipations, wins) property.

			Forms:
				/update event <event> points <points>
				/update event <event> enabled <enabled>
				/update event <event> effects <effects>
				/update user <user> points <points>
				/update user <user> partecipations <partecipations>
				/update user <user> wins <wins>
		*/
		// Check if the user is an bot-admin
		if !isAdmin(update.Message.From, utils) {
			// Respond and log with a message indicating that the user is not authorized to use this command
			SendUserNotAuthorizedMessage(update, data, utils)
			// Log the command failed execution
			FinalCommandLog("Unauthorized user", update, utils)
		} else {
			// Split the command arguments
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")
			// Check if the command arguments are in one of the above forms
			if len(cmdArgs) != 4 {
				// Respond with a message indicating that the command arguments are wrong
				cmdSyntax := "/update <\"event\"|\"user\"> <event|user> <\"points\"|\"enabled\"|\"effects\"|\"points\"|\"partecipations\"|\"wins\"> <points|enabled|effects|points|partecipations|wins>"
				SendWrongCommandSyntaxMessage(cmdSyntax, update, data, utils)
				// Log the command failed execution
				FinalCommandLog("Wrong command syntax", update, utils)
			} else {
				targetType := cmdArgs[0]
				switch targetType {
				case "event":
					// Get and check if the event exists
					eventKey := cmdArgs[1]
					if event, ok := events.Events.Map[eventKey]; !ok {
						// Respond with a message indicating that the event does not exist
						SendEntityNotFoundMessage("Evento", eventKey, update, data, utils)
						// Log the command failed execution
						FinalCommandLog("Event not found", update, utils)
					} else {
						targetProperty := cmdArgs[2]
						switch targetProperty {
						case "points":
							// Get and check if the points value is an integer number
							points, err := strconv.Atoi(cmdArgs[3])
							if err != nil {
								// Respond with a message indicating that the points value is not valid
								SendParameterNotValidMessage("points", "un numero intero", update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Update the Event.Points value
								events.Events.Map[eventKey] = &events.Event{Time: event.Time, Name: event.Name, Points: points, Enabled: event.Enabled, Effects: event.Effects, Activation: event.Activation, Partecipations: event.Partecipations}
								// Respond with command executed successfully
								SendPropertyUpdatedMessage("Event.Points", update, data, utils)
								// Log the /update command executed successfully
								FinalCommandLog("Event.Points updated", update, utils)
								SuccessResponseLog(update, utils)
							}
						case "enabled":
							// Get and check if the enabled value is a boolean
							enabled, err := strconv.ParseBool(cmdArgs[3])
							if err != nil {
								// Respond with a message indicating that the enabled value is not valid
								SendParameterNotValidMessage("enabled", "un booleano", update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Update the Event.Enabled value
								events.Events.Map[eventKey] = &events.Event{Time: event.Time, Name: event.Name, Points: event.Points, Enabled: enabled, Effects: event.Effects, Activation: event.Activation, Partecipations: event.Partecipations}
								// Respond with command executed successfully
								SendPropertyUpdatedMessage("Event.Enabled", update, data, utils)
								// Log the command executed successfully
								FinalCommandLog("Event.Enabled updated", update, utils)
								SuccessResponseLog(update, utils)
							}
						case "effects":
							// Get and check if the effects value is a slice of strings
							effectsNames, err := types.ParseSlice(cmdArgs[3])
							if err != nil {
								// Respond with a message indicating that the effects value is not valid
								SendParameterNotValidMessage("effects", "una lista di effetti validi", update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Get and check if the effects value is a slice of existing effects
								effects := make([]*structs.Effect, 0)
								wrongEffect := ""
								for _, effectName := range effectsNames {
									if _, ok := structs.Effects[effectName]; !ok {
										wrongEffect = effectName
										break
									}
									effects = append(effects, structs.Effects[effectName])
								}
								if wrongEffect == "" {
									// Update the Event.Effects value
									events.Events.Map[eventKey] = &events.Event{Time: event.Time, Name: event.Name, Points: event.Points, Enabled: event.Enabled, Effects: effects, Activation: event.Activation, Partecipations: event.Partecipations}
									// Respond with command executed successfully
									SendPropertyUpdatedMessage("Event.Effects", update, data, utils)
									// Log the command executed successfully
									FinalCommandLog("Event.Effects updated", update, utils)
									SuccessResponseLog(update, utils)
								} else {
									// Respond with a message indicating that the effect does not exist
									SendEntityNotFoundMessage("Effetto", wrongEffect, update, data, utils)
									// Log the command failed execution
									FinalCommandLog("Effect not found", update, utils)
								}
							}
						default:
							// Respond with a message indicating that the command arguments are wrong
							cmdSyntax := "/update <\"event\"|\"user\"> <event|user> <\"points\"|\"enabled\"|\"effects\"|\"points\"|\"partecipations\"|\"wins\"> <points|enabled|effects|points|partecipations|wins>"
							SendWrongCommandSyntaxMessage(cmdSyntax, update, data, utils)
							// Log the command failed execution
							FinalCommandLog("Wrong command syntax", update, utils)
						}
					}
				case "user":
					// Get and check if the user exists
					username := cmdArgs[1]
					var userKey int64
					for userID, user := range Users {
						if user != nil && user.UserName == username {
							userKey = userID
						}
					}
					if user, ok := Users[userKey]; !ok {
						// Respond with a message indicating that the user does not exist
						SendEntityNotFoundMessage("Utente", username, update, data, utils)
						// Log the command failed execution
						FinalCommandLog("User not found", update, utils)
					} else {
						targetProperty := cmdArgs[2]
						switch targetProperty {
						case "points":
							// Get and check if the points value is an integer number
							points, err := strconv.Atoi(cmdArgs[3])
							if err != nil {
								// Respond with a message indicating that the points value is not valid
								SendParameterNotValidMessage("points", "un numero intero", update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Update the User.TotalPoints value
								Users[userKey] = &structs.User{UserName: user.UserName, TotalPoints: points, TotalEventPartecipations: user.TotalEventPartecipations, TotalEventWins: user.TotalEventWins, TotalChampionshipPartecipations: user.TotalChampionshipPartecipations, TotalChampionshipWins: user.TotalChampionshipWins}
								// Respond with command executed successfully
								SendPropertyUpdatedMessage("TotalPoints", update, data, utils)
								// Log the command executed successfully
								FinalCommandLog("User.TotalPoints updated", update, utils)
								SuccessResponseLog(update, utils)
							}
						case "partecipations":
							// Get and check if the partecipations value is an integer number >= 0
							partecipations, err := strconv.Atoi(cmdArgs[3])
							if err != nil || partecipations < 0 {
								// Respond with a message indicating that the partecipations value is not valid
								SendParameterNotValidMessage("partecipations", "un numero intero positivo", update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Update the User.TotalEventPartecipations value
								Users[userKey] = &structs.User{UserName: user.UserName, TotalPoints: user.TotalPoints, TotalEventPartecipations: partecipations, TotalEventWins: user.TotalEventWins, TotalChampionshipPartecipations: user.TotalChampionshipPartecipations, TotalChampionshipWins: user.TotalChampionshipWins}
								// Respond with command executed successfully
								SendPropertyUpdatedMessage("TotalEventPartecipations", update, data, utils)
								// Log the command executed successfully
								FinalCommandLog("User.TotalEventPartecipations updated", update, utils)
								SuccessResponseLog(update, utils)
							}
						case "wins":
							// Get and check if the wins value is an integer number >= 0
							wins, err := strconv.Atoi(cmdArgs[3])
							if err != nil || wins < 0 {
								// Respond with a message indicating that the wins value is not valid
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Parametro <wins> deve essere un numero intero positivo.")
								SendMessage(msg, update, data, utils)
								// Log the command failed execution
								FinalCommandLog("Wrong command syntax", update, utils)
							} else {
								// Update the User.TotalEventWins value
								Users[userKey] = &structs.User{UserName: user.UserName, TotalPoints: user.TotalPoints, TotalEventPartecipations: user.TotalEventPartecipations, TotalEventWins: wins, TotalChampionshipPartecipations: user.TotalChampionshipPartecipations, TotalChampionshipWins: user.TotalChampionshipWins}
								// Respond with command executed successfully
								SendPropertyUpdatedMessage("TotalEventWins", update, data, utils)
								// Log the command executed successfully
								FinalCommandLog("User.TotalEventWins updated", update, utils)
								SuccessResponseLog(update, utils)
							}
						default:
							// Respond with a message indicating that the command arguments are wrong
							cmdSyntax := "/update <\"event\"|\"user\"> <event|user> <\"points\"|\"enabled\"|\"effects\"|\"points\"|\"partecipations\"|\"wins\"> <points|enabled|effects|points|partecipations|wins>"
							SendWrongCommandSyntaxMessage(cmdSyntax, update, data, utils)
							// Log the command failed execution
							FinalCommandLog("Wrong command syntax", update, utils)
						}
					}
				default:
					// Respond with a message indicating that the command arguments are wrong
					cmdSyntax := "/update <\"event\"|\"user\"> <event|user> <\"points\"|\"enabled\"|\"effects\"|\"points\"|\"partecipations\"|\"wins\"> <points|enabled|effects|points|partecipations|wins>"
					SendWrongCommandSyntaxMessage(cmdSyntax, update, data, utils)
					// Log the command failed execution
					FinalCommandLog("Wrong command syntax", update, utils)
				}
			}
		}
	}
}

// Check if a user is considered bot-admin (saved in .env file)
func isAdmin(user *tgbotapi.User, utils types.Utils) bool {
	adminUserIDStr := os.Getenv("TELEGRAM_ADMIN_ID")
	if adminUserIDStr == "" {
		utils.Logger.WithFields(logrus.Fields{
			"env": "TELEGRAM_ADMIN_ID",
		}).Panic("Env not set")
	}
	adminUserID, err := strconv.ParseInt(adminUserIDStr, 10, 64)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while parsing AdminUserIDStr")
	}

	return user.ID == adminUserID
}

// Send a message
func SendMessage(msg tgbotapi.MessageConfig, update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg.ReplyToMessageID = update.Message.MessageID
	message, err := data.Bot.Send(msg)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": message,
		}).Error("Error while sending message")
	}
}

func SendUserNotAuthorizedMessage(update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
	SendMessage(msg, update, data, utils)
}

func SendWrongCommandSyntaxMessage(cmdSyntax string, update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è: "+cmdSyntax)
	SendMessage(msg, update, data, utils)
}

func SendPropertyUpdatedMessage(property string, update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Proprietà '%v' aggiornata.", property))
	SendMessage(msg, update, data, utils)
}

func SendParameterNotValidMessage(parameter, expected string, update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Parametro <%v> non valido. Deve essere %v.", parameter, expected))
	SendMessage(msg, update, data, utils)
}

func SendEntityNotFoundMessage(entity string, entityValue any, update tgbotapi.Update, data types.Data, utils types.Utils) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v (%v) non trovato.", entity, entityValue))
	SendMessage(msg, update, data, utils)
}

func FinalCommandLog(msg string, update tgbotapi.Update, utils types.Utils) {
	utils.Logger.WithFields(logrus.Fields{
		"message": update.Message.Text,
		"sender":  update.Message.From.UserName,
		"chat":    update.Message.Chat.Title,
	}).Debug(msg)
}

func SuccessResponseLog(update tgbotapi.Update, utils types.Utils) {
	msg := fmt.Sprintf("Response to \"/%v\" command sent successfully", update.Message.Command())
	utils.Logger.Debug(msg)
}

func ComposeMessage(subMessages []string, args ...any) string {
	msg := ""
	for _, subMessage := range subMessages {
		msg += subMessage
	}
	return fmt.Sprintf(msg, args...)
}
