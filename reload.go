package main

import (
	"encoding/json"

	"github.com/MoraGames/clockyuwu/pkg/types"
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
