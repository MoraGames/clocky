package main

import (
	"fmt"
	"time"

	"github.com/MoraGames/clockyuwu/controller"
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func manageEvent(appUtils util.AppUtils, ctrler *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, event *model.Event) {
	if update.Message.Time().Format("15:04") == event.Time.Format("15:04") && update.Message.Text == event.Message {
		if event.Type != "second" {
			eventValid(appUtils, ctrler, bot, update, updateCurTime, event)
		} else if update.Message.Time().Format("15:04:05") == event.Time.Format("15:04:05") {
			eventValid(appUtils, ctrler, bot, update, updateCurTime, event)
		}
	}
	eventInvalid(appUtils, ctrler, bot, update, updateCurTime, event)
}

func eventValid(appUtils util.AppUtils, ctrler *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, event *model.Event) error {
	var user *model.User
	user, err := ctrler.GetUser(update.Message.From.ID)
	if err != nil {
		if err.Error() == "cannot get user not found" {
			err = ctrler.CreateUser(update.Message.From.ID, update.Message.From)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var msg tgbotapi.MessageConfig
	instance := event.Instances[len(event.Instances)-1]
	if !instance.Activated {
		//Win the event
		instance.Activated = true
		instance.ActivatedBy = user
		instance.ActivatedAt = updateCurTime
		instance.ArrivedAt = update.Message.Time()
		instance.Partecipations = append(instance.Partecipations, model.EventPartecipation{User: user, Time: updateCurTime})

		user.UserStats.TotalPoints += instance.Points
		user.UserStats.TotalEventPartecipations++
		user.UserStats.TotalEventWins++

		//Retrieve the partecipations data for the response message
		delay := updateCurTime.Sub(time.Date(instance.ArrivedAt.Year(), instance.ArrivedAt.Month(), instance.ArrivedAt.Day(), instance.ArrivedAt.Hour(), instance.ArrivedAt.Minute(), 0, 0, instance.ArrivedAt.Location()))

		//Prepare the message
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Complimenti %v! +%v punti per te.\nHai impiegato +%vs.", update.Message.From.UserName, instance.Points, delay.Seconds()))
	} else {
		//Check if the user has already partecipated to the event
		alreadyPartecipated := false
		for _, partecipant := range instance.Partecipations {
			if partecipant.User.TelegramUser.ID == user.TelegramUser.ID {
				alreadyPartecipated = true
				break
			}
		}

		if alreadyPartecipated {
			//Already partecipated to the event

			//Prepare the message
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v hai già partecipato a questo evento.", update.Message.From.UserName))
		} else {
			//Partecipate (without win) to the event
			instance.Partecipations = append(instance.Partecipations, model.EventPartecipation{User: user, Time: updateCurTime})

			user.UserStats.TotalEventPartecipations++

			//Retrieve the partecipations data for the response message
			delay := updateCurTime.Sub(time.Date(instance.ArrivedAt.Year(), instance.ArrivedAt.Month(), instance.ArrivedAt.Day(), instance.ArrivedAt.Hour(), instance.ArrivedAt.Minute(), 0, 0, instance.ArrivedAt.Location()))
			delta := updateCurTime.Sub(instance.ActivatedAt)

			//Prepare the message
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs.", instance.ActivatedBy, delta.Seconds(), delay.Seconds()))
		}
	}
	event.Instances[len(event.Instances)-1] = instance

	//Update the user
	err = ctrler.UpdateUser(user.TelegramUser.ID, user)
	if err != nil {
		return err
	}

	//Update the event
	err = ctrler.UpdateEvent(event.Message, event)
	if err != nil {
		return err
	}

	//Send the message
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	return nil
}

func eventInvalid(appUtils util.AppUtils, ctrler *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, event *model.Event) error {
	var msg tgbotapi.MessageConfig
	if update.Message.Time().Format("15:04") == event.Time.Add(-1*time.Minute).Format("15:04") {
		//Prepare the message
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Troppo veloce %v, non è ancora il momento di partecipare per le %v.", update.Message.From.UserName, event.Message))
	} else if update.Message.Time().Format("15:04") == event.Time.Add(1*time.Minute).Format("15:04") {
		//Prepare the message
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v sei arrivato troppo tardi, le %v non potevano restare ad aspettarti oltre.", update.Message.From.UserName, event.Message))
	} else {
		return nil
	}

	//Send the message
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
	return nil
}
