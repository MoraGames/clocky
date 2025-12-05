package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func manageMigrations() {

}

func AddTelegramUserToExistingUser(msgFrom *tgbotapi.User) {
	if user, exist := Users[msgFrom.ID]; exist {
		user.TelegramUser = msgFrom
		Users[msgFrom.ID] = user
	}
}
