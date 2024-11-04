package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	EventMutex sync.Mutex
	UserMutex  sync.Mutex
)

func manageEvent(eventUpdate utils.EventUpdate) error {
	//Retrieve the user data
	user, err := getUserData(eventUpdate.Update.Message.From)
	if err != nil {
		return err
	}

	//Get the last instance of the event (the current one)
	var instance *model.EventInstance
	if instance, err = App.Controller.GetLastEventInstanceByEvent(eventUpdate.Event); err != nil {
		return err
	}

	var msg tgbotapi.MessageConfig
	//Check if the user has won/lost the event
	EventMutex.Lock()
	if !instance.Activated {
		//Win the event
		instance.Activated = true
		instance.ActivatedBy = user
		instance.ActivatedAt = eventUpdate.Update.Message.Time()
		instance.ArrivedAt = eventUpdate.BotTime
		App.Controller.CreateEventPartecipationForEventInstance(&model.EventPartecipation{
			EventInstance: instance,
			User:          user,
			//WIP: Need to be converted to model.Chat
			//Chat:          eventUpdate.Update.Message.Chat,
			Time: eventUpdate.BotTime,
		})
		EventMutex.Unlock()

		//Update the user stats
		App.Controller.UpdateUserStats(user, instance.Points, true)

		//Retrieve the partecipations data for the response message
		delay := eventUpdate.BotTime.Sub(time.Date(instance.ArrivedAt.Year(), instance.ArrivedAt.Month(), instance.ArrivedAt.Day(), instance.ArrivedAt.Hour(), instance.ArrivedAt.Minute(), 0, 0, instance.ArrivedAt.Location()))

		//Prepare the message
		var msgGreetings string
		switch {
		case instance.Points < 0:
			msgGreetings = "Accidenti"
		case instance.Points == 0:
			msgGreetings = "Peccato"
		case instance.Points > 0:
			msgGreetings = "Complimenti"
		}
		msg = tgbotapi.NewMessage(eventUpdate.Update.Message.Chat.ID, fmt.Sprintf("%v %v! +%v punti per te.\nHai impiegato +%vs.", msgGreetings, eventUpdate.Update.Message.From.UserName, instance.Points, delay.Seconds()))
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
			msg = tgbotapi.NewMessage(eventUpdate.Update.Message.Chat.ID, fmt.Sprintf("%v hai già partecipato a questo evento.", eventUpdate.Update.Message.From.UserName))
		} else {
			//Partecipate (without win) to the event
			App.Controller.CreateEventPartecipationForEventInstance(&model.EventPartecipation{
				EventInstance: instance,
				User:          user,
				//WIP: Need to be converted to model.Chat
				//Chat:          eventUpdate.Update.Message.Chat,
				Time: eventUpdate.BotTime,
			})

			//Update the user stats
			App.Controller.UpdateUserStats(user, 0, false)

			//Retrieve the partecipations data for the response message
			delay := eventUpdate.BotTime.Sub(time.Date(instance.ArrivedAt.Year(), instance.ArrivedAt.Month(), instance.ArrivedAt.Day(), instance.ArrivedAt.Hour(), instance.ArrivedAt.Minute(), 0, 0, instance.ArrivedAt.Location()))
			delta := eventUpdate.BotTime.Sub(instance.ActivatedAt)

			//Prepare the message
			msg = tgbotapi.NewMessage(eventUpdate.Update.Message.Chat.ID, fmt.Sprintf("L'evento è già stato attivato da %v +%vs fa.\nHai impiegato +%vs.", instance.ActivatedBy, delta.Seconds(), delay.Seconds()))
		}
	}

	//WIP: Commented the old-reworked code to mantain a trace of the checks and logic to implement.
	//		This part of the code is not working in any version and should not be never used.
	//		This part of the code need to be removed meanwhile the new working code is implemented.
	//TODO: Implements the full events management logic starting from the old code.
	/*
		if !instance.Activated {
			//Win the event
			instance.Activated = true
			instance.ActivatedBy = user
			instance.ActivatedAt = updateCurTime
			instance.ArrivedAt = update.Message.Time()
			instance.Partecipations = append(instance.Partecipations, &model.EventPartecipation{User: user, Time: updateCurTime})

			user.UserStats.TotalPoints += instance.Points
			user.UserStats.TotalEventsPartecipations++
			user.UserStats.TotalEventsWins++

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

				user.UserStats.TotalEventsPartecipations++

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

	*/

	//Send the message
	msg.ReplyToMessageID = eventUpdate.Update.Message.MessageID
	App.BotAPI.Send(msg)

	return nil
}

func getUserData(telegramUser *tgbotapi.User) (*model.User, error) {
	//Lock the mutex to prevent multiple creation of the same user. Defer the unlock
	UserMutex.Lock()
	defer UserMutex.Unlock()

	//Check if the user exists, if exist retrieve it's data
	user, err := App.Controller.GetUserByTelegramID(telegramUser.ID)
	if err != nil {
		//If the user does not exist, create it
		if err.Error() == "cannot get user not found" {
			if err = App.Controller.CreateUser(telegramUser); err != nil {
				return nil, err
			}
			//Retrieve the user data
			user, err = App.Controller.GetUserByTelegramID(telegramUser.ID)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}

/*
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
*/
