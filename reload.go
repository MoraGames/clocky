package main

import (
	"encoding/json"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/sirupsen/logrus"
)

func reloadStatus(reloads []types.Reload, utils types.Utils) {
	utils.Logger.Info("Reloading data from files")

	numOfFail, numOfFailFunc, numOfOkay, numOfOkayFunc := 0, 0, 0, 0
	for _, reload := range reloads {
		hasFailed := false

		utils.Logger.WithFields(logrus.Fields{
			"IfFail()": reload.IfFail != nil,
			"IfOkay()": reload.IfOkay != nil,
		}).Debug("Reloading " + reload.FileName)

		file, err := App.FilesRoot.ReadFile(reload.FileName)
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

// If reloaded events or championship are expired, reset them to default values
func ResetExpiredData(utils types.Utils) {
	if events.CurrentChampionship.Expiration.Before(time.Now()) {
		utils.Logger.WithFields(logrus.Fields{
			"exp": events.CurrentChampionship.Expiration,
			"now": time.Now(),
		}).Info("Resetting championship to default values due to expiration date passed")
		events.CurrentChampionship.Reset(
			structs.GetRanking(Users, structs.RankScopeChampionship, true),
			&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
			types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
		)
	}
	if events.Events.Expiration.Before(time.Now()) {
		utils.Logger.WithFields(logrus.Fields{
			"exp": events.Events.Expiration,
			"now": time.Now(),
		}).Info("Resetting events to default values due to expiration date passed")
		events.Events.Reset(true, &types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1}, utils)
	}
}
