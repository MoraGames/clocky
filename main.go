package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
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

	logFile, err := os.OpenFile("files/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	l.SetOutput(mw)

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

	//get current time location
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Time location not get (using UTC)")
	}

	//set the default chat ID
	defChatIDstr := os.Getenv("TELEGRAM_DEFAULT_CHAT_ID")
	if defChatIDstr == "" {
		l.WithFields(logrus.Fields{
			"env": "TELEGRAM_DEFAULT_CHAT_ID",
		}).Warn("Env not set")
	}

	defChatID, err := strconv.ParseInt(defChatIDstr, 10, 64)
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Error while parsing TELEGRAM_DEFAULT_CHAT_ID to int64")
	}

	//set the gocron events reset
	gcScheduler := gocron.NewScheduler(timeLocation)
	gcJob, err := gcScheduler.Every(1).Day().At("23:58").Do(
		events.Events.Reset, defChatIDstr != "", &types.WriteMessageData{Bot: bot, ChatID: defChatID, ReplyMessageID: -1}, types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"})
	if err != nil {
		l.WithFields(logrus.Fields{
			"gcJob": gcJob,
			"error": err,
		}).Error("GoCron job not set")
	}

	updates := bot.GetUpdatesChan(u)
	l.WithFields(logrus.Fields{
		"debugMode": bot.Debug,
		"timeout":   u.Timeout,
	}).Debug("Update channel retreived")

	// Read from specified files and reload the data into the structs
	ReloadStatus(
		[]types.Reload{
			{FileName: "files/sets.json", DataStruct: &events.Sets},
			{FileName: "files/events.json", DataStruct: &events.Events},
			{FileName: "files/users.json", DataStruct: &Users},
		},
		types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"},
	)

	gcScheduler.StartAsync()
	run(types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: bot, Updates: updates})
	gcScheduler.Stop()
}

func ReloadStatus(reloads []types.Reload, utils types.Utils) {
	utils.Logger.WithFields(logrus.Fields{
		"reloads": reloads,
	}).Debug("Reloading status")

	for _, reload := range reloads {
		file, err := os.ReadFile(reload.FileName)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
				"err":  err,
			}).Error("Error while reading file")
		}

		if len(file) != 0 {
			err = json.Unmarshal(file, reload.DataStruct)
			if err != nil {
				utils.Logger.WithFields(logrus.Fields{
					"data": reload.DataStruct,
					"err":  err,
				}).Error("Error while unmarshalling data")
			}
		} else {
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Warn("File is empty")
		}
	}
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
