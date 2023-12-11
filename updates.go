package main

import (
	"time"

	"github.com/MoraGames/clockyuwu/controller"
	"github.com/MoraGames/clockyuwu/pkg/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func manageUpdates(appUtils util.AppUtils, ctrler *controller.Controller, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		//log the update received
		updateCurTime := time.Now()
		appUtils.Logger.WithFields(logrus.Fields{
			"chat":    update.FromChat().Title,
			"curTime": updateCurTime.Format(appUtils.TimeFormat),
		}).Info("Update received")

		//check if the update is a callback query
		if update.CallbackQuery != nil {
			//log the callback query received
			appUtils.Logger.WithFields(logrus.Fields{}).Info("CallbackQuery received")

			//TODO: Manage CallbackQuery
		}

		//check if the update is a message
		if update.Message != nil {
			//log the message received
			appUtils.Logger.WithFields(logrus.Fields{
				"message": update.Message.Text,
				"sender":  update.Message.From.UserName,
				"msgTime": update.Message.Time().Format(appUtils.TimeFormat),
				"curTime": updateCurTime.Format(appUtils.TimeFormat),
			}).Info("Message received")

			//check if the message is a command
			if update.Message.IsCommand() {
				manageCommand(appUtils, ctrler, bot, update, updateCurTime)
				continue
			}

			//check if the message is an event
			if test, eventType, err := ctrler.IsEvent(update.Message.Text); err != nil {
				appUtils.Logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error while checking if the message is an event")
			} else if test {
				manageEvent(appUtils, ctrler, bot, update, updateCurTime, eventType)
				continue
			}
		}
	}

	return nil
}
