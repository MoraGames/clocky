package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/internal/app"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/pkg/utils"
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func DefineDefaultCronJobs() {
	// Daily reset cronjob - at 23:59:30
	if _, err := App.GocronScheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23, 59, 30))),
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
		gocron.WithEventListeners(
			// Daily summary report
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				rankingString := ""
				ranking := structs.GetRanking(Users, structs.RankScopeDay, true)
				for i, rankEntry := range ranking {
					rankingString += fmt.Sprintf("%d] %s: %d\n", i+1, rankEntry.Username, rankEntry.Points)
				}
				entities, text := utils.ParseToEntities(ComposeMessage(
					[]string{
						"__**Ecco la classifica di giornata:**__\n\n",
						rankingString + "\n",
						"Un giro di applausi per tutti i partecipanti di oggi, ma adesso preparatevi tutti, un nuovo giorno di sfide sta già per sorgere!",
					},
				), TelegramUsersList)
				respMsg := tgbotapi.NewMessage(App.DefaultChatID, text)
				respMsg.Entities = entities
				App.BotAPI.Send(respMsg)
			}),
		),
	); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"job": "DailyResetCronjob",
			"err": err,
		}).Error("GoCron job not set")
	}

	// Championship reset cronjob - every 2 Sundays at 23:59:25
	if _, err := App.GocronScheduler.NewJob(
		gocron.WeeklyJob(2, gocron.NewWeekdays(time.Sunday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 25))),
		gocron.NewTask(func() {
			events.CurrentChampionship.Reset(
				structs.GetRanking(Users, structs.RankScopeChampionship, true),
				&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
				types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
			)

			// Reward the users based on their performance
			// Then reset the championship user's stats
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

	// Championship 7 days reminder cronjob - every 2 Sundays at 23:59:25
	if _, err := App.GocronScheduler.NewJob(
		gocron.WeeklyJob(2, gocron.NewWeekdays(time.Sunday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 25))),
		gocron.NewTask(func() {
			entities, text := utils.ParseToEntities(ComposeMessage(
				[]string{
					"📯 __**A tutti i giocatori e giocatrici:**__\n",
					"Siamo già a metà dell'opera! __Rimangono 7 giorni alla fine del campionato__.\n",
					"Ricordiamo tutti che una volta proclamato il vincitore, __egli avrà diritto esclusivo ad un bonus__ per l'intera durata del campionato successivo!\n",
					"Affrettatevi dunque, nulla ancora è perduto!",
				},
			), TelegramUsersList)
			respMsg := tgbotapi.NewMessage(App.DefaultChatID, text)
			respMsg.Entities = entities
			App.BotAPI.Send(respMsg)
		}),
		gocron.WithName("Championship7DaysReminderCronjob"),
		gocron.JobOption(gocron.WithStartDateTimePast(
			utils.NextInWeekdayAtTime(time.Now(), time.Sunday, 23, 59, 25).AddDate(0, 0, -7),
		)),
	); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"job": "Championship7DaysReminderCronjob",
			"err": err,
		}).Error("GoCron job not set")
	}

	// Championship 1 day reminder cronjob - every 2 Saturdays at 23:59:25
	if _, err := App.GocronScheduler.NewJob(
		gocron.WeeklyJob(2, gocron.NewWeekdays(time.Saturday), gocron.NewAtTimes(gocron.NewAtTime(23, 59, 25))),
		gocron.NewTask(func() {
			entities, text := utils.ParseToEntities(ComposeMessage(
				[]string{
					"🔥 __**Attenzione giocatori!**__\n",
					"Mancano esattamente 24 ore alla fine del Campionato in corso!\n",
					"Questa è la vostra ultima occasione per __scalare la classifica e conquistare il titolo di **Clocky Champion**__!\n",
					"Buon divertimento a tutti e che vinca il migliore! 🏆\n\n",
					"P.S.:\n",
					"//Per chi non dovesse riuscirci, non preoccupatevi: presto tutti i punteggi verranno azzerati e la sfida ricomincerà da capo.//\n",
				},
			), TelegramUsersList)
			respMsg := tgbotapi.NewMessage(App.DefaultChatID, text)
			respMsg.Entities = entities
			App.BotAPI.Send(respMsg)
		}),
		gocron.WithName("Championship1DayReminderCronjob"),
	); err != nil {
		App.Logger.WithFields(logrus.Fields{
			"job": "Championship1DayReminderCronjob",
			"err": err,
		}).Error("GoCron job not set")
	}

	App.Logger.WithFields(logrus.Fields{
		"gcJobs": utils.StringifyJobs(App.GocronScheduler.Jobs()),
	}).Info("GoCron jobs set")
}

func ChampionshipUserRewardAndReset(users map[int64]*structs.User, writeMsgData *types.WriteMessageData, utilsVar types.Utils) {
	// Reward the user that have won the championship
	ranking := structs.GetRanking(Users, structs.RankScopeChampionship, true)
	for userId := range Users {
		if user, ok := Users[userId]; ok && user != nil {
			// Remove the reigning leader and reigning podium effects
			Users[userId].RemoveEffect(structs.ReigningLeader)

			// Check if the user is the winner of the championship
			if len(ranking) > 0 && ranking[0].UserTelegramID == userId {
				// Update the data structure of deserving users
				Users[userId].TotalChampionshipWins++
				Users[userId].AddEffect(structs.ReigningLeader)
			}
			// Check if the user is in the top 3 of the ranking
			for r := 0; r < 3 && r < len(ranking); r++ {
				if ranking[r].UserTelegramID == userId {
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
	todayRewardedUsers := make([]events.DailyRewardedUser, 0)
	for userId := range Users {
		if user, ok := Users[userId]; ok && user != nil {
			// Check if the user has participated in at least 10% of the enabled events of the day
			if Users[userId].DailyEventPartecipations >= int(math.Round(float64(dailyEnabledEvents)*0.10)) {
				// Update the data structure of deserving users
				Users[userId].DailyPartecipationStreak++

				// Reward the user (hints)
				choosenSets := ManageDailyRewardMessage(userId, writeMsgData, utilsVar)
				todayRewardedUsers = append(todayRewardedUsers, events.DailyRewardedUser{User: user.Minimize(), Sets: choosenSets})

				// Check if the user has won at least 25% of the events in which he participated
				if Users[userId].DailyEventWins >= int(math.Round(float64(Users[userId].DailyEventPartecipations)*0.25)) {
					Users[userId].DailyActivityStreak++
				} else {
					Users[userId].DailyActivityStreak = 0
				}
			} else {
				Users[userId].DailyPartecipationStreak = 0
			}

			// Reward the user (activity streak)
			Users[userId].RemoveEffect(structs.NoNegative)
			Users[userId].RemoveEffect(structs.ConsistentParticipant1)
			Users[userId].RemoveEffect(structs.ConsistentParticipant2)
			if Users[userId].DailyActivityStreak >= 7 {
				Users[userId].AddEffect(structs.ConsistentParticipant1)
			}
			if Users[userId].DailyActivityStreak >= 14 {
				Users[userId].AddEffect(structs.ConsistentParticipant2)
			}
			if Users[userId].DailyActivityStreak >= 21 {
				Users[userId].AddEffect(structs.NoNegative)
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
	var finalPositionMessage, effectRewardMessage, effectAppliedMessage string
	switch rankPosition {
	case 0:
		finalPositionMessage = "You are the new Clocky Champion!"
		effectRewardMessage = "For this you are rewarded with a special bonus effects for the entire duration of the next championship (if you choose to participate in it)"
		effectAppliedMessage = "\n\nYou are the proud owner of the effect \"Reigning Leader\", wich grants you a +1 points bonus to every event you will win.\n Congratulations again!"
	case 1:
		finalPositionMessage = "You're standing on the 2nd step of the podium."
		effectRewardMessage = "Unfortunately, you were so close to victory, but that wasn't enough to obtain any special effects. We're sure we'll see you with victory in your grasp next season!"
	case 2:
		finalPositionMessage = "You're standing on the 3rd step of the podium."
		effectRewardMessage = "Unfortunately, you were so close to victory, but that wasn't enough to obtain any special effects. We're sure we'll see you with victory in your grasp next season!"
	}
	text := fmt.Sprintf("Congratulations %v!\n%v %v%v", Users[userId].UserName, finalPositionMessage, effectRewardMessage, effectAppliedMessage)

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

func ManageDailyRewardMessage(userId int64, writeMsgData *types.WriteMessageData, utils types.Utils) []string {
	// Generate the reward message informations
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var availableSets, choosedSets []string
	var setsEvents [][]*events.Event
	availableSets = append(availableSets, events.Events.Stats.EnabledSets...)
	for (len(choosedSets) < 3 && len(choosedSets) < len(availableSets)) || (len(setsEvents) < 20 && len(choosedSets) < len(availableSets)) {
		choosedSets = append(choosedSets, availableSets[r.Intn(events.Events.Stats.EnabledSetsNum-len(choosedSets))])
		setsEvents = append(setsEvents, events.EventsOf(events.SetsFunctions[choosedSets[len(choosedSets)-1]]))
	}
	var numEffects []int
	for _, setEvents := range setsEvents {
		for _, event := range setEvents {
			numEffects = append(numEffects, len(event.Effects))
		}
	}

	// Generate the reward message
	text := fmt.Sprintf("Congratulations %v!\nYou have won %v events and for this you are rewarded with an hint for the new day.\nHere are some of the events and relative effects that are surely active in the next 24 hours:\n\n", Users[userId].UserName, Users[userId].DailyEventWins)
	for i, setEvents := range setsEvents {
		text += fmt.Sprintf("Events of the Set %q (%v events with %v effects):\n", choosedSets[i], len(setEvents), numEffects[i])
		for _, event := range setEvents {
			text += fmt.Sprintf(" | %q", event.Name)
			eventEffects := event.StringifyEffects()
			if eventEffects != "[]" {
				text += fmt.Sprintf("  with %v", eventEffects)
			}
			text += "\n"
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

	return choosedSets
}

// This function was supposed to be written in events/championship.go, but due to dependencies with App variable, it's here instead.
func UpdateChampionshipCronjobs(utilsVar types.Utils) {
	// Check if the scheduler is initialized
	if App.GocronScheduler == nil {
		App.Logger.Error("Scheduler not initialized before reload")
		return
	}

	// Find the job IDs
	jobIDs := map[string]uuid.UUID{
		"ChampionshipResetCronjob":         uuid.Nil,
		"Championship7DaysReminderCronjob": uuid.Nil,
		"Championship1DayReminderCronjob":  uuid.Nil,
	}
	for _, job := range App.GocronScheduler.Jobs() {
		if _, exist := jobIDs[job.Name()]; exist {
			jobIDs[job.Name()] = job.ID()
		}
	}
	for name, id := range jobIDs {
		if id == uuid.Nil {
			App.Logger.WithFields(logrus.Fields{
				"job": name,
			}).Error("GoCron job not found during reload")
		}
	}
	// Update the cronjobs
	championshipStartDate := events.CurrentChampionship.StartDate.In(time.Local)

	for name, id := range jobIDs {
		if id != uuid.Nil {
			switch name {
			case "ChampionshipResetCronjob":
				if _, err := App.GocronScheduler.Update(
					id,
					gocron.WeeklyJob(2, gocron.NewWeekdays(championshipStartDate.Weekday()), gocron.NewAtTimes(gocron.NewAtTime(uint(championshipStartDate.Hour()), uint(championshipStartDate.Minute()), uint(championshipStartDate.Second())))),
					gocron.NewTask(func() {
						events.CurrentChampionship.Reset(
							structs.GetRanking(Users, structs.RankScopeChampionship, true),
							&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
							types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
						)

						// Reward the users based on their performance
						// Then reset the championship user's stats
						ChampionshipUserRewardAndReset(
							Users,
							&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
							types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
						)
					}),
					gocron.WithName("ChampionshipResetCronjob"),
					gocron.WithStartAt(gocron.WithStartDateTimePast(
						events.CurrentChampionship.StartDate.In(time.Local),
					)),
				); err != nil {
					App.Logger.WithFields(logrus.Fields{
						"job": "ChampionshipResetCronjob",
						"err": err,
					}).Error("GoCron job not updated during reload")
				}
			case "Championship7DaysReminderCronjob":
				if _, err := App.GocronScheduler.Update(
					id,
					gocron.WeeklyJob(2, gocron.NewWeekdays(championshipStartDate.Weekday()), gocron.NewAtTimes(gocron.NewAtTime(uint(championshipStartDate.Hour()), uint(championshipStartDate.Minute()), uint(championshipStartDate.Second())))),
					gocron.NewTask(func() {
						entities, text := utils.ParseToEntities(ComposeMessage(
							[]string{
								"📯 __**A tutti i giocatori e giocatrici:**__\n",
								"Siamo già a metà dell'opera! __Rimangono 7 giorni alla fine del campionato__.\n",
								"Ricordiamo tutti che una volta proclamato il vincitore, __egli avrà diritto esclusivo ad un bonus__ per l'intera durata del campionato successivo!\n",
								"Affrettatevi dunque, nulla ancora è perduto!",
							},
							app.Name,
						), TelegramUsersList)
						respMsg := tgbotapi.NewMessage(App.DefaultChatID, text)
						respMsg.Entities = entities
						App.BotAPI.Send(respMsg)
					}),
					gocron.WithName("Championship7DaysReminderCronjob"),
					gocron.WithStartAt(gocron.WithStartDateTimePast(
						events.CurrentChampionship.StartDate.In(time.Local).AddDate(0, 0, -7),
					)),
				); err != nil {
					App.Logger.WithFields(logrus.Fields{
						"job": "Championship7DaysReminderCronjob",
						"err": err,
					}).Error("GoCron job not updated during reload")
				}
			case "Championship1DayReminderCronjob":
				if _, err := App.GocronScheduler.Update(
					id,
					gocron.WeeklyJob(2, gocron.NewWeekdays(championshipStartDate.AddDate(0, 0, -1).Weekday()), gocron.NewAtTimes(gocron.NewAtTime(uint(championshipStartDate.Hour()), uint(championshipStartDate.Minute()), uint(championshipStartDate.Second())))),
					gocron.NewTask(func() {
						entities, text := utils.ParseToEntities(ComposeMessage(
							[]string{
								"🔥 __**Attenzione giocatori!**__\n",
								"Mancano esattamente 24 ore alla fine del Campionato in corso!\n",
								"Questa è la vostra ultima occasione per __scalare la classifica e conquistare il titolo di **Clocky Champion**__!\n",
								"Buon divertimento a tutti e che vinca il migliore! 🏆\n\n",
								"P.S.:\n",
								"//Per chi non dovesse riuscirci, non preoccupatevi: presto tutti i punteggi verranno azzerati e la sfida ricomincerà da capo.//\n",
							},
							app.Name,
						), TelegramUsersList)
						respMsg := tgbotapi.NewMessage(App.DefaultChatID, text)
						respMsg.Entities = entities
						App.BotAPI.Send(respMsg)
					}),
					gocron.WithName("Championship1DayReminderCronjob"),
					gocron.WithStartAt(gocron.WithStartDateTimePast(
						events.CurrentChampionship.StartDate.In(time.Local).AddDate(0, 0, -1),
					)),
				); err != nil {
					App.Logger.WithFields(logrus.Fields{
						"job": "Championship1DayReminderCronjob",
						"err": err,
					}).Error("GoCron job not updated during reload")
				}
			}
		}
	}
}
