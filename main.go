package main

import (
	"log"
	"os"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	"github.com/MoraGames/clockyuwu/pkg/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	//get the configurations
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//setup the logger
	l := logger.NewLogger(conf.Log.Type, conf.Log.Format, conf.Log.Level)
	l.WithFields(logrus.Fields{
		"lvl": conf.Log.Level,
		"typ": conf.Log.Type,
		"frm": conf.Log.Format,
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

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	run(types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: bot, Updates: updates})
}
