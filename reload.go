package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
			} else if reload.Validate != nil && !reload.Validate(utils) {
				hasFailed = true
				backupExpiredFile(reload.FileName, utils)
				utils.Logger.WithFields(logrus.Fields{
					"file": reload.FileName,
				}).Warn("Reloaded data is expired")
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

func backupExpiredFile(fileName string, utils types.Utils) {
	content, err := os.ReadFile(filepath.Join("files", fileName))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{"file": fileName, "err": err}).Error("Unable to read expired file for backup")
		return
	}

	backupDir := filepath.Join("files", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		utils.Logger.WithFields(logrus.Fields{"file": fileName, "err": err}).Error("Unable to create expired files backup directory")
		return
	}

	baseName := strings.ReplaceAll(filepath.Base(fileName), ".", "_")
	backupName := fmt.Sprintf("%s.%s.json", baseName, time.Now().Format("20060102_150405.000000000"))
	backupPath := filepath.Join(backupDir, backupName)
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		utils.Logger.WithFields(logrus.Fields{"file": fileName, "backup": backupPath, "err": err}).Error("Unable to write expired file backup")
		return
	}
	utils.Logger.WithFields(logrus.Fields{"file": fileName, "backup": backupPath}).Info("Expired file backed up")
}
