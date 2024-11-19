package main

import (
	"encoding/json"
	"fmt"
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
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	//get the configurations
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//setup the logger
	l := logger.NewLogger(
		logger.LoggerOutput{
			LogWriter:     logger.StringToWriter(conf.Logger.Console.Writer),
			LogType:       conf.Logger.Console.Type,
			LogTimeFormat: conf.Logger.Console.TimeFormat,
			LogLevel:      conf.Logger.Console.Level,
		},
		logger.LoggerOutput{
			LogWriter: &lumberjack.Logger{
				Filename: conf.Logger.File.Location,
				MaxSize:  conf.Logger.File.MaxSize, // MB
			},
			LogType:       conf.Logger.File.Type,
			LogTimeFormat: conf.Logger.File.TimeFormat,
			LogLevel:      conf.Logger.File.Level,
		},
	)
	l.WithFields(logrus.Fields{
		"typ": conf.Logger.Console.Type,
		"lvl": conf.Logger.Console.Level,
		"fmt": conf.Logger.Console.TimeFormat,
	}).Debug("Output ", conf.Logger.Console.Writer, " set")
	l.WithFields(logrus.Fields{
		"typ": conf.Logger.File.Type,
		"lvl": conf.Logger.File.Level,
		"fmt": conf.Logger.File.TimeFormat,
	}).Debug("Output ", conf.Logger.File.Location, " set")
	l.WithFields(logrus.Fields{
		"outs": []string{conf.Logger.Console.Writer, conf.Logger.File.Location},
	}).Info("Logger initialized")

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

	//get current time location
	timeLocation, err := time.LoadLocation("Local")
	if err != nil {
		l.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Time location not get (using UTC)")
	}

	//set the gocron events reset
	gcScheduler := gocron.NewScheduler(timeLocation)
	gcJob, err := gcScheduler.Every(1).Day().At("23:59").Do(
		func() {
			// Get the number of enabled events for the ended day
			dailyEnabledEvents := events.Events.Stats.EnabledEventsNum

			// Reset the events
			events.Events.Reset(
				true,
				&types.WriteMessageData{Bot: bot, ChatID: defChatID, ReplyMessageID: -1},
				types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"},
			)

			// Reward the users where DailyEventWins >= 30% of TotalDailyEventWins
			// Then reset the daily user's stats (unconditionally)
			DailyUserRewardAndReset(
				Users,
				dailyEnabledEvents,
				&types.WriteMessageData{Bot: bot, ChatID: defChatID, ReplyMessageID: -1},
				types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"},
			)
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
	}).Info("Update channel retreived")

	// Read from specified files and reload the data into the structs
	ReloadStatus(
		[]types.Reload{
			{FileName: "files/sets.json", DataStruct: &events.SetsJson, IfOkay: events.AssignSetsFromSetsJson, IfFail: events.AssignSetsWithDefault},
			{FileName: "files/events.json", DataStruct: &events.Events, IfOkay: nil, IfFail: events.AssignEventsWithDefault},
			{FileName: "files/users.json", DataStruct: &Users, IfOkay: nil, IfFail: nil},
			{FileName: "files/pinnedMessage.json", DataStruct: &events.PinnedResetMessage, IfOkay: nil, IfFail: nil},
			{FileName: "files/hints.json", DataStruct: &events.HintRewardedUsers, IfOkay: nil, IfFail: nil},
		},
		types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"},
	)

	gcScheduler.StartAsync()
	run(types.Utils{Config: conf, Logger: l, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: bot, Updates: updates})
	gcScheduler.Stop()
}

func ReloadStatus(reloads []types.Reload, utils types.Utils) {
	utils.Logger.Info("Reloading data from files")

	numOfFail, numOfFailFunc, numOfOkay, numOfOkayFunc := 0, 0, 0, 0
	for _, reload := range reloads {
		hasFailed := false

		utils.Logger.WithFields(logrus.Fields{
			"IfFail()": reload.IfFail != nil,
			"IfOkay()": reload.IfOkay != nil,
		}).Debug("Reloading " + reload.FileName)

		file, err := os.ReadFile(reload.FileName)
		if err != nil {
			hasFailed = true
			utils.Logger.WithFields(logrus.Fields{
				"file": reload.FileName,
				"err":  err,
			}).Error("Error while reading file")
		} else if len(file) != 0 {
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
				}).Debug("Reload.IfFail() executed")
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
				}).Debug("Reload.IfOkay() executed")
			}
		}
	}

	utils.Logger.WithFields(logrus.Fields{
		"fails":     numOfFail,
		"failsFunc": numOfFailFunc,
		"okays":     numOfOkay,
		"okaysFunc": numOfOkayFunc,
		"total":     len(reloads),
	}).Info("Reloading data completed")
}

func DailyUserRewardAndReset(users map[int64]*structs.User, dailyEnabledEvents int, writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Reward the users where DailyEventWins >= 30% of TotalDailyEventWins
	// Then reset the daily user's stats (unconditionally)
	todayRewardedUsers := make([]*structs.UserMinimal, 0)
	for userId := range Users {
		if user, ok := Users[userId]; ok && user != nil {
			// Check if the user has participated in at least 15% of the enabled events of the day and if he has won at least 25% of the events in which he participated
			if Users[userId].DailyEventPartecipations >= int(math.Round(float64(dailyEnabledEvents)*0.15)) && Users[userId].DailyEventWins >= int(math.Round(float64(Users[userId].DailyEventPartecipations)*0.25)) {
				// Update the data structure of deserving users
				todayRewardedUsers = append(todayRewardedUsers, user.Minimize())

				// Reward the user
				ManageRewardMessage(userId, writeMsgData, utils)
			}

			// Reset the daily user's stats
			Users[userId].DailyPoints = 0
			Users[userId].DailyEventPartecipations = 0
			Users[userId].DailyEventWins = 0
		}
	}

	// Update UserHintMessages
	events.HintRewardedUsers[time.Now().Format("02-01-2006")] = todayRewardedUsers

	// Save the users
	file, err := json.MarshalIndent(Users, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to marshal Users data",
		}).Error("Error while marshalling data")
		utils.Logger.Error(Users)
	}
	err = os.WriteFile("files/users.json", file, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to write Users data",
		}).Error("Error while writing data")
		utils.Logger.Error(Users)
	}

	// Save the hints sent
	file, err = json.MarshalIndent(events.HintRewardedUsers, "", " ")
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to marshal HintRewards data",
		}).Error("Error while marshalling data")
		utils.Logger.Error(events.HintRewardedUsers)
	}
	err = os.WriteFile("files/hints.json", file, 0644)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to write HintRewards data",
		}).Error("Error while writing data")
		utils.Logger.Error(events.HintRewardedUsers)
	}
}

func ManageRewardMessage(userId int64, writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Generate the reward message informations
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomSet := events.Events.Stats.EnabledSets[r.Intn(events.Events.Stats.EnabledSetsNum)]
	setEvents := events.EventsOf(events.SetsFunctions[randomSet])
	numEffects := 0
	for _, event := range setEvents {
		numEffects += len(event.Effects)
	}

	// Generate the reward message
	text := fmt.Sprintf("Congratulations %v!\nYou have won %v/%v events you entered and for this you are rewarded with an hint for the new day.\nHere are some of the events and relative effects that are surely active in the next 24 hours:\n\nEvents of the Set %q (%v events with %v effects):\n", Users[userId].UserName, Users[userId].DailyEventWins, Users[userId].DailyEventPartecipations, randomSet, len(setEvents), numEffects)
	for _, event := range setEvents {
		text += fmt.Sprintf(" | %q", event.Name)
		eventEffects := event.StringifyEffects()
		if eventEffects != "[]" {
			text += fmt.Sprintf("  with %v", eventEffects)
		}
		text += "\n"
	}

	// Send the reward message
	msg := tgbotapi.NewMessage(userId, text)
	message, err := writeMsgData.Bot.Send(msg)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": message,
		}).Error("Error while sending message")
	}
}

func WriteMessage(bot *tgbotapi.BotAPI, chatID int64, replyMessageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMessageID != -1 {
		msg.ReplyToMessageID = replyMessageID
	}
	bot.Send(msg)
}
