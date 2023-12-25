package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/go-co-op/gocron"
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

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	l.SetOutput(mw)

	//get current time location
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Time location not get (using UTC)")
	}

	//set the gocron events reset
	gcScheduler := gocron.NewScheduler(timeLocation)
	gcJob, err := gcScheduler.Every(1).Day().At("23:58").Do(events.Events.Reset, false, types.WriteMessageData{})
	if err != nil {
		l.WithFields(logrus.Fields{
			"gcJob": gcJob,
			"error": err,
		}).Error("GoCron job not set")
	}

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

	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 180

	updates := bot.GetUpdatesChan(u)
	l.WithFields(logrus.Fields{
		"debugMode": bot.Debug,
		"timeout":   u.Timeout,
	}).Debug("Update channel retreived")

	gcScheduler.StartAsync()
	run(types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: bot, Updates: updates})
	gcScheduler.Stop()
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
