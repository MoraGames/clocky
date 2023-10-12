package main

import (
	"os"
	"log"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Utils struct {
	conf *config.Config
	log *logrus.Logger
}

type Data struct {
	bot *tgbotapi.BotAPI
	updates *tgbotapi.UpdatesChannel
}

func main() {
	//get the configurations
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//setup the logger
	l := logger.NewLogger(conf.Log.Level, conf.Log.Type)
	l.WithFields(logrus.Fields{
		"lvl": conf.Log.Level,
		"typ": conf.Log.Type,
	}).Debug("Logger initialized")

	//link Telegram API
	apiToken := os.Getenv("TelegramAPIToken")
	if apiToken == "" {
		l.WithFields(logrus.Fields{
			"env": "TelegramAPIToken",
		}).Panic("Env not set")
	}

	//get the bot API
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while getting bot API")
	}

	l.WithFields(logrus.Fields{
		"id":       bot.Self.ID,
		"username": bot.Self.UserName,
	}).Info("Account authorized")
		
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	run(l, bot, updates)
}
