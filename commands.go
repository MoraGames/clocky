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

func manageCommands(update tgbotapi.Update, utils types.Utils, data types.Data, curTime time.Time, eventKey events.EventKey) {
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
	case "reset":
		if !isAdmin(update.Message.From, utils) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
			msg.ReplyToMessageID = update.Message.MessageID
			data.Bot.Send(msg)
			utils.Logger.WithFields(logrus.Fields{
				"usr": update.Message.From.UserName,
				"cmd": update.Message.Command(),
			}).Debug("Unauthorized user")
		} else {
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")
			if len(cmdArgs) != 1 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /reset <events|users>")
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"msg": update.Message.Text,
				}).Debug("Wrong command")
			} else {
				switch cmdArgs[0] {
				case "events":
					events.Events.Reset()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Eventi resettati")
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.Debug("Events resetted")
				case "users":
					Users = make(map[int64]*structs.User)
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

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Utenti resettati")
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.Debug("Users resetted")
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /reset <events|users>")
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.WithFields(logrus.Fields{
						"usr": update.Message.From.UserName,
						"cmd": update.Message.Command(),
					}).Debug("Wrong command")
				}
			}
		}
	case "update":
		if !isAdmin(update.Message.From, utils) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
			msg.ReplyToMessageID = update.Message.MessageID
			data.Bot.Send(msg)
			utils.Logger.WithFields(logrus.Fields{
				"usr": update.Message.From.UserName,
				"cmd": update.Message.Command(),
			}).Debug("Unauthorized user")
		} else {
			cmdArgs := strings.Split(update.Message.CommandArguments(), " ")
			if len(cmdArgs) != 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /update <event> <points>")
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"msg": update.Message.Text,
				}).Debug("Wrong command")
			} else {
				eventKeyString := cmdArgs[0]
				if event, ok := events.Events[events.EventKey(eventKeyString)]; !ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Evento non trovato")
					msg.ReplyToMessageID = update.Message.MessageID
					data.Bot.Send(msg)
					utils.Logger.WithFields(logrus.Fields{
						"usr": update.Message.From.UserName,
						"msg": update.Message.Text,
					}).Debug("Event not found")
				} else {
					points, err := strconv.Atoi(cmdArgs[1])
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "parametro points deve essere un numero.")
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)
						utils.Logger.WithFields(logrus.Fields{
							"usr": update.Message.From.UserName,
							"msg": update.Message.Text,
						}).Debug("Wrong command")
					} else {
						events.Events[events.EventKey(eventKeyString)] = &events.EventValue{Points: points, Activated: event.Activated, ActivatedBy: event.ActivatedBy, ActivatedAt: event.ActivatedAt, ArrivedAt: event.ArrivedAt, Partecipations: event.Partecipations}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Evento aggiornato")
						msg.ReplyToMessageID = update.Message.MessageID
						data.Bot.Send(msg)
						utils.Logger.Debug("Event updated")
					}
				}
			}
		}
	case "status":
		cmdArgs := strings.Split(update.Message.CommandArguments(), " ")
		if len(cmdArgs) == 1 && cmdArgs[0] == "get" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "bot is currently: "+data.Status)
			msg.ReplyToMessageID = update.Message.MessageID
			data.Bot.Send(msg)
			utils.Logger.Debug("Status getted")
		} else if len(cmdArgs) == 2 && cmdArgs[0] == "set" {
			if !isAdmin(update.Message.From, utils) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad usare questo comando")
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"cmd": update.Message.Command(),
				}).Debug("Unauthorized user")
			} else if cmdArgs[1] == "online" {
				data.Status = "online"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "bot is currently: "+data.Status)
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"sts": data.Status,
				}).Info("Bot status set to online")
			} else if cmdArgs[1] == "offline" {
				data.Status = "offline"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "bot is currently: "+data.Status)
				msg.ReplyToMessageID = update.Message.MessageID
				data.Bot.Send(msg)
				utils.Logger.WithFields(logrus.Fields{
					"usr": update.Message.From.UserName,
					"sts": data.Status,
				}).Info("Bot status set to online")
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Il comando è /status <get|set <online|offline>> ")
			msg.ReplyToMessageID = update.Message.MessageID
			data.Bot.Send(msg)
			utils.Logger.WithFields(logrus.Fields{
				"usr": update.Message.From.UserName,
				"msg": update.Message.Text,
			}).Debug("Wrong command")
		}
	}
}

func isAdmin(user *tgbotapi.User, utils types.Utils) bool {
	adminUserIDStr := os.Getenv("AdminUserID")
	if adminUserIDStr == "" {
		utils.Logger.WithFields(logrus.Fields{
			"env": "AdminUserID",
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
