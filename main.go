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
	"github.com/MoraGames/clockyuwu/pkg/utils"
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var App utils.Application

func init() {
	//get the configurations
	var err error
	App.Config, err = config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

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
	if _, err = App.GocronScheduler.NewJob(
		gocron.DailyJob(2, gocron.NewAtTimes(gocron.NewAtTime(23, 59, 50))),
		gocron.NewTask(func() {
			// Get the number of enabled events for the ended day
			dailyEnabledEvents := events.Events.Stats.EnabledEventsNum

			// Reset the events
			events.Events.Reset(
				true,
				&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
				types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
			)

			// Reward the users based on their performance
			// Then reset the daily user's stats (unconditionally)
			DailyUserRewardAndReset(
				Users,
				dailyEnabledEvents,
				&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
				types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
			)
		}),
		gocron.WithName("DailyResetCronjob"),
	); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"job": "DailyResetCronjob",
			"err": err,
		}).Error("GoCron job not set")
	}
	if _, err = App.GocronScheduler.NewJob(
		gocron.WeeklyJob(2, gocron.NewWeekdays(time.Sunday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 40))),
		gocron.NewTask(func() {
			ChampionshipUserRewardAndReset(
				Users,
				&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
				types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
			)
		}),
		gocron.WithName("ChampionshipResetCronjob"),
		//gocron.WithStartAt(gocron.WithStartDateTimePast()) //Add by reload if successful
	); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"job": "ChampionshipResetCronjob",
			"err": err,
		}).Error("GoCron job not set")
	}
	App.Logger.WithFields(logrus.Fields{
		"gcJobs": utils.StringifyJobs(App.GocronScheduler.Jobs()),
	}).Info("GoCron jobs set")
}

func main() {
	reloadStatus(
		[]types.Reload{
			{FileName: "files/sets.json", DataStruct: &events.SetsJson, IfOkay: events.AssignSetsFromSetsJson, IfFail: events.AssignSetsWithDefault},
			{FileName: "files/events.json", DataStruct: &events.Events, IfOkay: nil, IfFail: events.AssignEventsWithDefault},
			{FileName: "files/users.json", DataStruct: &Users, IfOkay: nil, IfFail: nil},
			{FileName: "files/pinnedMessage.json", DataStruct: &events.PinnedResetMessage, IfOkay: nil, IfFail: nil},
			{FileName: "files/hints.json", DataStruct: &events.HintRewardedUsers, IfOkay: nil, IfFail: nil},
			{FileName: "files/championship.json", DataStruct: &events.CurrentChampionship, IfOkay: UpdateChampionshipResetCronjob, IfFail: events.AssignChampionshipWithDefault},
		},
		types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: "15:04:05.000000 MST -07:00"},
	)
	manageMigrations()

	App.GocronScheduler.Start()
	run(types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: "15:04:05.000000 MST -07:00"}, types.Data{Bot: App.BotAPI, Updates: App.Updates})
	App.GocronScheduler.Shutdown()
}

func ChampionshipUserRewardAndReset(users map[int64]*structs.User, writeMsgData *types.WriteMessageData, utilsVar types.Utils) {
	// Reward the user that have won the championship
	ranking := utils.GetRanking(Users, utils.RankScopeChampionship)
	for userId := range Users {
		if user, ok := Users[userId]; ok && user != nil {
			// Remove the reigning leader and reigning podium effects
			Users[userId].RemoveEffect(structs.ReigningLeader)
			Users[userId].RemoveEffect(structs.ReigningPodium)

			// Check if the user is in the top 3 of the ranking
			for r := 0; r < 3 && r < len(ranking); r++ {
				if ranking[r].UserTelegramID == userId {
					// Update the data structure of deserving users
					if r == 0 {
						Users[userId].TotalChampionshipWins++
						Users[userId].AddEffect(structs.ReigningLeader)
					} else {
						Users[userId].AddEffect(structs.ReigningPodium)
					}

					// Reward the user
					ManageChampionshipRewardMessage(userId, r, writeMsgData, utilsVar)
				}
			}

			// Reset and update the championship user's stats
			if Users[userId].ChampionshipEventPartecipations > 0 && Users[userId].ChampionshipEventWins > 0 {
				Users[userId].TotalChampionshipPartecipations++
			}
			Users[userId].ChampionshipPoints = 0
			Users[userId].ChampionshipEventPartecipations = 0
			Users[userId].ChampionshipEventWins = 0
		}
	}

	// Save the users
	file, err := json.MarshalIndent(Users, "", " ")
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to marshal Users data",
		}).Error("Error while marshalling data")
		utilsVar.Logger.Error(Users)
	}
	err = os.WriteFile("files/users.json", file, 0644)
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to write Users data",
		}).Error("Error while writing data")
		utilsVar.Logger.Error(Users)
	}
}

func DailyUserRewardAndReset(users map[int64]*structs.User, dailyEnabledEvents int, writeMsgData *types.WriteMessageData, utilsVar types.Utils) {
	// Reward the users where DailyEventWins >= 30% of TotalDailyEventWins
	// Then reset the daily user's stats (unconditionally)
	todayRewardedUsers := make([]*structs.UserMinimal, 0)
	for userId := range Users {
		if user, ok := Users[userId]; ok && user != nil {
			// Check if the user has participated in at least 15% of the enabled events of the day
			if Users[userId].DailyEventPartecipations >= int(math.Round(float64(dailyEnabledEvents)*0.15)) {
				// Update the data structure of deserving users
				Users[userId].DailyPartecipationStreak++
			} else {
				Users[userId].DailyPartecipationStreak = 0
			}

			// Check if the user has participated in at least 15% of the enabled events of the day and if he has won at least 25% of the events in which he participated
			if Users[userId].DailyEventPartecipations >= int(math.Round(float64(dailyEnabledEvents)*0.15)) && Users[userId].DailyEventWins >= int(math.Round(float64(Users[userId].DailyEventPartecipations)*0.25)) {
				// Update the data structure of deserving users
				todayRewardedUsers = append(todayRewardedUsers, user.Minimize())
				Users[userId].DailyActivityStreak++

				// Reward the user
				ManageDailyRewardMessage(userId, writeMsgData, utilsVar)
			} else {
				Users[userId].DailyActivityStreak = 0
			}

			Users[userId].RemoveEffect(structs.NoNegative)
			Users[userId].RemoveEffect(structs.ConsistentParticipant1)
			Users[userId].RemoveEffect(structs.ConsistentParticipant2)

			// Check if the user have an active activity streak of at least 7/21 days
			if Users[userId].DailyActivityStreak >= 21 {
				Users[userId].AddEffect(structs.NoNegative)
				Users[userId].AddEffect(structs.ConsistentParticipant2)
			} else if Users[userId].DailyActivityStreak >= 7 {
				Users[userId].AddEffect(structs.ConsistentParticipant1)
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
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to marshal Users data",
		}).Error("Error while marshalling data")
		utilsVar.Logger.Error(Users)
	}
	err = os.WriteFile("files/users.json", file, 0644)
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to write Users data",
		}).Error("Error while writing data")
		utilsVar.Logger.Error(Users)
	}

	// Save the hints sent
	file, err = json.MarshalIndent(events.HintRewardedUsers, "", " ")
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to marshal HintRewards data",
		}).Error("Error while marshalling data")
		utilsVar.Logger.Error(events.HintRewardedUsers)
	}
	err = os.WriteFile("files/hints.json", file, 0644)
	if err != nil {
		utilsVar.Logger.WithFields(logrus.Fields{
			"err": err,
			"msg": "Unable to write HintRewards data",
		}).Error("Error while writing data")
		utilsVar.Logger.Error(events.HintRewardedUsers)
	}
}

func ManageChampionshipRewardMessage(userId int64, rankPosition int, writeMsgData *types.WriteMessageData, utils types.Utils) {
	// Generate the reward message
	var finalPositionMessage, effectAppliedMessage string
	switch rankPosition {
	case 0:
		finalPositionMessage = "You are the new Clocky Champion! For"
		effectAppliedMessage = "You are the proud owner of the effect \"Reigning Leader\", wich grants you a +3 points bonus to every event you will participate.\n Congratulations again!"
	case 1:
		finalPositionMessage = "You're standing on the 2nd step of the podium, and for"
		effectAppliedMessage = "You are an owner of the effect \"Reigning Podium\", wich grants you a +2 points bonus to every event you will participate."
	case 2:
		finalPositionMessage = "You're standing on the 3rd step of the podium, and for"
		effectAppliedMessage = "You are an owner of the effect \"Reigning Podium\", wich grants you a +2 points bonus to every event you will participate."
	}
	text := fmt.Sprintf("Congratulations %v!\n%v this you are rewarded with a special bonus effect for the entire duration of the next championship (if you choose to participate in it):\n\n%v", Users[userId].UserName, finalPositionMessage, effectAppliedMessage)

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

func ManageDailyRewardMessage(userId int64, writeMsgData *types.WriteMessageData, utils types.Utils) {
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

// This function was supposed to be written in events/championship.go, but due to dependencies with App variable, it's here instead.
func UpdateChampionshipResetCronjob(utils types.Utils) {
	// Check if the scheduler is initialized
	if App.GocronScheduler == nil {
		App.Logger.Error("Scheduler not initialized before reload")
		return
	}

	// Find the job
	var found bool
	var jobID uuid.UUID
	for _, job := range App.GocronScheduler.Jobs() {
		if job.Name() == "ChampionshipResetCronjob" {
			found = true
			jobID = job.ID()
			break
		}
	}
	if !found {
		App.Logger.Error("ChampionshipResetCronjob not found in scheduler")
		return
	}

	// Update the cronjob
	job, err := App.GocronScheduler.Update(
		jobID,
		gocron.WeeklyJob(2, gocron.NewWeekdays(time.Sunday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 40))),
		gocron.NewTask(func() {
			ChampionshipUserRewardAndReset(
				Users,
				&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
				types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
			)
		}),
		gocron.WithName("ChampionshipResetCronjob"),
		gocron.WithStartAt(gocron.WithStartDateTimePast(
			events.CurrentChampionship.StartDate,
		)),
	)

	jobNextRun, err := job.NextRun()
	if err != nil {
		App.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Error while getting next run time of ChampionshipResetCronjob")
		return
	}
	App.Logger.WithFields(logrus.Fields{
		"name": events.CurrentChampionship.Name,
		"date": events.CurrentChampionship.StartDate,
		"next": jobNextRun,
	}).Info("Championship schedule restored from reload")
}
