package main

import (
	"log"
	"os"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/controller"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	"github.com/MoraGames/clockyuwu/pkg/util"
	"github.com/MoraGames/clockyuwu/repo/mock"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

//TODO: Do an E/R diagram and change all models before integrating the planetScaleDB repository.
//TODO: Update all mock repositories to use the new models. Implements planetScaleDB repository.
//TODO: Refactor all controller module.

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
	}).Debug("Logger initialized")

	//link Telegram API
	apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if apiToken == "" {
		l.WithFields(logrus.Fields{
			"env": "TELEGRAM_API_TOKEN",
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

	//setup the updates channel
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120

	updates := bot.GetUpdatesChan(u)
	l.WithFields(logrus.Fields{
		"updates": updates,
	}).Debug("Updates channel initialized")

	//initialize the appUtils data struct
	appUtils := util.AppUtils{
		ConfigApp:  conf.App,
		Logger:     l,
		TimeFormat: "15:04:05.000000 MST -07:00",
	}

	//initialize the controller data struct
	controller := controller.NewController(
		mock.NewBonusRepo(),
		mock.NewChampionshipRepo(),
		mock.NewChatRepo(),
		mock.NewEventRepo(),
		mock.NewPartecipationRepo(),
		mock.NewRecordRepo(),
		mock.NewSetRepo(),
		mock.NewUserRepo(),
		l,
	)

	//run the bot over the updates channel
	err = manageUpdates(appUtils, controller, bot, updates)
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while managing updates")
	}
}
