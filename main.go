package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
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
	gcJob, err := gcScheduler.Every(1).Day().At("23:59").Do(
		func() {
			// Get the number of enabled events for the ended day
			DailyEnabledEvents := events.Events.Stats.EnabledEventsNum

			// Reset the events
			events.Events.Reset(
				true,
				&types.WriteMessageData{Bot: bot, ChatID: defChatID, ReplyMessageID: -1},
				types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"},
			)

			// Reward the users where DailyEventWins >= 30% of TotalDailyEventWins
			// Then reset the daily user's stats (unconditionally)
			for userId := range Users {
				if user, ok := Users[userId]; ok && user != nil {
					// Check if the user has participated in at least 20% of the enabled events of the day and if he has won at least 30% of the events in which he participated
					if Users[userId].DailyEventPartecipations >= int(math.Round(float64(DailyEnabledEvents)*0.2)) && Users[userId].DailyEventWins >= int(math.Round(float64(Users[userId].DailyEventPartecipations)*0.3)) {
						r := rand.New(rand.NewSource(time.Now().UnixNano()))
						randomSet := events.Events.Stats.EnabledSets[r.Intn(events.Events.Stats.EnabledSetsNum)]
						setEvents := events.EventsOf(events.SetsFunctions[randomSet])

						text := fmt.Sprintf("Congratulations! You have won %v/%v events you entered and for this you are rewarded with an hint for the new day.\nHere are some of the events surely active in the next 24 hours:\n\nEvents of the Set %q:", Users[userId].DailyEventWins, Users[userId].DailyEventPartecipations, randomSet)
						for _, event := range setEvents {
							text += fmt.Sprintf(" | %q\n", event)
						}

						msg := tgbotapi.NewMessage(userId, text)
						message, error := bot.Send(msg)
						if error != nil {
							l.WithFields(logrus.Fields{
								"err": error,
								"msg": message,
							}).Error("Error while sending message")
						}
					}

					Users[userId].DailyPoints = 0
					Users[userId].DailyEventPartecipations = 0
					Users[userId].DailyEventWins = 0
				}
			}

			// Save the users
			file, err := json.MarshalIndent(Users, "", " ")
			if err != nil {
				l.WithFields(logrus.Fields{
					"err":  err,
					"note": "preoccupati",
				}).Error("Error while marshalling data")
				l.Error(Users)
			}
			err = os.WriteFile("files/users.json", file, 0644)
			if err != nil {
				l.WithFields(logrus.Fields{
					"err":  err,
					"note": "preoccupati tanto",
				}).Error("Error while writing data")
				l.Error(Users)
			}
		},
	)
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
			{FileName: "files/sets.json", DataStruct: &events.SetsJson, IfOkay: events.AssignSetsFromSetsJson, IfFail: events.AssignSetsWithDefault},
			{FileName: "files/events.json", DataStruct: &events.Events, IfOkay: nil, IfFail: events.AssignEventsWithDefault},
			{FileName: "files/users.json", DataStruct: &Users, IfOkay: nil, IfFail: nil},
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
	}).Info("Reloading data")

	numOfFail, numOfFailFunc, numOfOkay, numOfOkayFunc := 0, 0, 0, 0
	for _, reload := range reloads {
		hasFailed := false

		utils.Logger.WithFields(logrus.Fields{
			"IfFail() exist": reload.IfFail != nil,
			"IfOkay() exist": reload.IfOkay != nil,
		}).Debug("Reloading " + reload.FileName)

		file, err := os.ReadFile(reload.FileName)
		if err != nil {
			hasFailed = true
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
				"err":  err,
			}).Error("Error while reading file")
		}

		if len(file) != 0 {
			err = json.Unmarshal(file, reload.DataStruct)
			if err != nil {
				hasFailed = true
				utils.Logger.WithFields(logrus.Fields{
					"data": reload.DataStruct,
					"err":  err,
				}).Error("Error while unmarshalling data")
			}
		} else {
			hasFailed = true
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Error("File is empty")
		}

		if hasFailed {
			numOfFail++

			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Warn("Reloading has failed")

			if reload.IfFail != nil {
				numOfFailFunc++
				reload.IfFail(utils)
				utils.Logger.WithFields(logrus.Fields{
					"file": reload.FileName,
				}).Info("Reload.IfFail() executed")
			}
		} else {
			numOfOkay++
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
			}).Debug("Reloading has succeed")

			if reload.IfOkay != nil {
				numOfOkayFunc++
				reload.IfOkay(utils)
				utils.Logger.WithFields(logrus.Fields{
					"file": reload.FileName,
				}).Info("Reload.IfOkay() executed")
			}
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"fails":     numOfFail,
		"failsFunc": numOfFailFunc,
		"okays":     numOfOkay,
		"okaysFunc": numOfOkayFunc,
		"totat":     len(reloads),
	}).Info("Reloading data completed")
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
