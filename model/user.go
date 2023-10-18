package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type User struct {
	TelegramUser        *tgbotapi.User
	Points              int
	EventPartecipations int
	EventWins           int
}
