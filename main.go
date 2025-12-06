package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MoraGames/clockyuwu/config"
	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/internal/app"
	"github.com/MoraGames/clockyuwu/pkg/logger"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var App app.Application
var envModeFlag string

func init() {
	//define the flags and their aliases
	flag.StringVar(&envModeFlag, "envmode", "", "Select the environment to use (matches .env.<envmode>)")
	flag.StringVar(&envModeFlag, "env", "", "Alias of \"envmode\"")
}

func main() {
	//get the configurations
	flag.Parse()
	App.EnvMode = config.ResolveEnvMode(envModeFlag)

	var err error
	App.Config, err = config.NewConfig(App.EnvMode)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Configurations loaded for environment:", App.EnvMode)

	//setup the logger
	App.Logger = logger.NewLogger(
		logger.LoggerOutput{
			LogWriter:     logger.StringToWriter(App.Config.Logger.Console.Writer),
			LogType:       App.Config.Logger.Console.Type,
			LogTimeFormat: App.Config.Logger.Console.TimeFormat,
			LogLevel:      App.Config.Logger.Console.Level,
		},
		logger.LoggerOutput{
			LogWriter: &lumberjack.Logger{
				Filename: App.Config.Logger.File.Location,
				MaxSize:  App.Config.Logger.File.MaxSize, // MB
			},
			LogType:       App.Config.Logger.File.Type,
			LogTimeFormat: App.Config.Logger.File.TimeFormat,
			LogLevel:      App.Config.Logger.File.Level,
		},
	)
	App.Logger.WithFields(logrus.Fields{
		"typ": App.Config.Logger.Console.Type,
		"lvl": App.Config.Logger.Console.Level,
		"fmt": App.Config.Logger.Console.TimeFormat,
	}).Debug("Output ", App.Config.Logger.Console.Writer, " set")
	App.Logger.WithFields(logrus.Fields{
		"typ": App.Config.Logger.File.Type,
		"lvl": App.Config.Logger.File.Level,
		"fmt": App.Config.Logger.File.TimeFormat,
	}).Debug("Output ", App.Config.Logger.File.Location, " set")
	App.Logger.WithFields(logrus.Fields{
		"outs": []string{App.Config.Logger.Console.Writer, App.Config.Logger.File.Location},
	}).Info("Logger initialized")

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

	App.BotAPI.Debug = false
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 180

	//get the updates channel
	App.Updates = App.BotAPI.GetUpdatesChan(upd)
	App.Logger.WithFields(logrus.Fields{
		"debugMode": App.BotAPI.Debug,
		"timeout":   upd.Timeout,
	}).Info("Update channel retreived")

	defaultChatEnv := os.Getenv("TELEGRAM_DEFAULT_CHAT_ID")
	if defaultChatEnv == "" {
		App.Logger.WithFields(logrus.Fields{
			"env": "TELEGRAM_DEFAULT_CHAT_ID",
		}).Warn("Env not set")
	}

	App.DefaultChatID, err = strconv.ParseInt(defaultChatEnv, 10, 64)
	if err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Error while parsing TELEGRAM_DEFAULT_CHAT_ID to int64")
	}

	//get current time location
	if _, err = time.LoadLocation("Local"); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Time location not get (using UTC)")
	}

	App.TimeFormat = "15:04:05.000000 MST -07:00"

	//create the gocron scheduler
	App.GocronScheduler, err = gocron.NewScheduler()
	if err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error while creating GoCron scheduler")
	}

	//define the default cron jobs for the application scheduler
	DefineDefaultCronJobs()

	//try to reload the status from files
	reloadStatus(
		[]types.Reload{
			{FileName: "files/sets.json", DataStruct: &events.SetsJson, IfOkay: events.AssignSetsFromSetsJson, IfFail: events.AssignSetsWithDefault},
			{FileName: "files/events.json", DataStruct: &events.Events, IfOkay: nil, IfFail: events.AssignEventsWithDefault},
			{FileName: "files/users.json", DataStruct: &Users, IfOkay: nil, IfFail: nil},
			{FileName: "files/pinnedMessage.json", DataStruct: &events.PinnedResetMessage, IfOkay: nil, IfFail: nil},
			{FileName: "files/hints.json", DataStruct: &events.HintRewardedUsers, IfOkay: nil, IfFail: nil},
			{FileName: "files/championship.json", DataStruct: &events.CurrentChampionship, IfOkay: UpdateChampionshipResetCronjob, IfFail: events.AssignChampionshipWithDefault},
			{FileName: "files/pinnedChampionshipMessage.json", DataStruct: &structs.PinnedChampionshipResetMessage, IfOkay: nil, IfFail: nil},
		},
		types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: "15:04:05.000000 MST -07:00"},
	)
	GenerateTelegramUsersList()

	//manage data migrations
	manageMigrations()

	//start the scheduler and run the bot
	App.GocronScheduler.Start()
	run(types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: App.BotAPI, Updates: App.Updates})
	App.GocronScheduler.Shutdown()
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
