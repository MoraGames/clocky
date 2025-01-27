package utils

import (
	"github.com/MoraGames/clockyuwu/config"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config          *config.Config
	Logger          *logrus.Logger
	BotAPI          *tgbotapi.BotAPI
	Updates         tgbotapi.UpdatesChannel
	GocronScheduler *gocron.Scheduler
	TimeFormat      string
}
