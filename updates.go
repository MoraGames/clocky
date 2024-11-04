package main

import (
	"time"

	"github.com/MoraGames/clockyuwu/pkg/utils"
	"github.com/sirupsen/logrus"
)

func manageUpdates() error {
	for upd := range App.Updates {
		//log the update received
		botTime := time.Now()
		App.Logger.WithFields(logrus.Fields{
			"chat":    upd.FromChat().Title,
			"botTime": botTime.Format(App.TimesFormat),
		}).Debug("Update received")

		//check if the update is a callback query
		if upd.CallbackQuery != nil {
			//log the callback query received
			App.Logger.WithFields(logrus.Fields{}).Info("CallbackQuery received")

			go manageCallbackQuery(utils.CallbackQueryUpdate{Update: upd, BotTime: botTime})
			continue //TODO: this could be wrong (test what happens with callback queries and how to manage them)
		}

		//check if the update is a message
		if upd.Message != nil {
			/*
				//log the message received
				App.Logger.WithFields(logrus.Fields{
					"sender":  update.Message.From.UserName,
					"msgTime": update.Message.Time().Format(appUtils.TimeFormat),
					"botTime": updateCurTime.Format(appUtils.TimeFormat),
				}).Debug("Message received")
			*/

			//check if the message is a command
			if upd.Message.IsCommand() {
				//log the command received
				App.Logger.WithFields(logrus.Fields{
					"command": upd.Message.Command(),
					"botTime": botTime.Format(App.TimesFormat),
				}).Info("Command received")

				go manageCommand(utils.CommandUpdate{Update: upd, BotTime: botTime})
				continue
			}

			//check if the message is an active event sent in the correct time
			if test, event := App.Controller.IsEvent(upd.Message.Text); test && event.Enabled && event.Time.Format("15:04") == upd.Message.Time().Format("15:04") {
				//log the event received
				App.Logger.WithFields(logrus.Fields{
					"event":   upd.Message.Text,
					"botTime": botTime.Format(App.TimesFormat),
				}).Info("Event received")

				//WIP: Since the function is now asynchronous, the function need to manage a common mutex to avoid syncronization problems.
				go manageEvent(utils.EventUpdate{Update: upd, BotTime: botTime, Event: event})
				continue
			}
		}
	}

	return nil
}
