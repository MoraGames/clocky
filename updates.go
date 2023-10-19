package main

import (
	"time"

	"github.com/MoraGames/clockyuwu/pkg/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func manageUpdate(appUtils util.AppUtils, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		}
	}
}
