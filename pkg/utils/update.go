package utils

import (
	"time"

	"github.com/MoraGames/clockyuwu/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Update interface {
		GetUpdate() tgbotapi.Update
		GetBotTime() time.Time
	}
	CallbackQueryUpdate struct {
		Update  tgbotapi.Update
		BotTime time.Time
	}
	CommandUpdate struct {
		Update  tgbotapi.Update
		BotTime time.Time
	}
	EventUpdate struct {
		Update  tgbotapi.Update
		BotTime time.Time
		Event   *model.Event
	}
)

func (cu *CallbackQueryUpdate) GetUpdate() tgbotapi.Update {
	return cu.Update
}
func (cu *CallbackQueryUpdate) GetBotTime() time.Time {
	return cu.BotTime
}

func (cu *CommandUpdate) GetUpdate() tgbotapi.Update {
	return cu.Update
}
func (cu *CommandUpdate) GetBotTime() time.Time {
	return cu.BotTime
}

func (eu *EventUpdate) GetUpdate() tgbotapi.Update {
	return eu.Update
}
func (eu *EventUpdate) GetBotTime() time.Time {
	return eu.BotTime
}
