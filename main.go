package main

import (
	"log"
	"os"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	"github.com/MoraGames/clockyuwu/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

//TODO: Do an E/R diagram and change all models before integrating the planetScaleDB repository.
//TODO: Update all mock repositories to use the new models. Implements planetScaleDB repository.
//TODO: Refactor all controller module.
//TODO: PlanetScaleDB is now a payed service. Find and implement a new cloud db service for free.
//		- Oracles Cloud Free Tier: https://www.oracle.com/it/cloud/free/
//		- Neo4j-Graph AuraDB: https://neo4j.com/pricing/
//		- DataStax AstraDB: https://www.datastax.com/pricing/astra-db
//		Other suggestions are Google Spreadsheet or Airtable (both are sheets-like db).

var App utils.Application

func init() {
	//get the configurations
	var err error
	App.Config, err = config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//setup the logger
	App.Logger = logger.NewLogger(App.Config.Logger.Type, App.Config.Logger.Format, App.Config.Logger.Level, App.Config.Logger.Rotation)
	App.Logger.WithFields(logrus.Fields{
		"lvl": App.Config.Logger.Level,
		"rot": App.Config.Logger.Rotation,
	}).Debug("Logger initialized")

	//link Telegram API
	apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	if apiToken == "" {
		App.Logger.WithFields(logrus.Fields{
			"env": "TELEGRAM_API_TOKEN",
		}).Panic("Env not set")
	}

	//get the bot API
	App.BotAPI, err = tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while getting bot API")
	}

	App.Logger.WithFields(logrus.Fields{
		"id":       App.BotAPI.Self.ID,
		"username": App.BotAPI.Self.UserName,
	}).Info("Account authorized")

	//setup the updates channel
	App.BotAPI.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120

	App.Updates = App.BotAPI.GetUpdatesChan(u)
	App.Logger.Debug("Updates channel initialized")

	//define other application infos
	App.TimesFormat = "15:04:05.000000 MST -07:00"
	App.Author = "@MoraGames"
	App.Name = "ClockyMaster"
	App.Version = "2.0.0 RC-1"

	/*
		//initialize the controller data struct
		App.Controller = controller.NewController(
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
	*/
}

func main() {
	//execute the migrations function
	if err := manageMigrations(); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while managing migrations")
	}

	//run the bot over the updates channel
	if err := manageUpdates(); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Error while managing updates")
	}
}
