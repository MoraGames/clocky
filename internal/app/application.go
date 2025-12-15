package app

import (
	"os"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	Name    = "Clocky"
	Version = "0.4.3"
)

type Application struct {
	EnvMode         string
	Config          *config.Config
	Logger          *logrus.Logger
	BotAPI          *tgbotapi.BotAPI
	FilesRoot       *os.Root
	DefaultChatID   int64
	GocronScheduler gocron.Scheduler
	TimeFormat      string
}
