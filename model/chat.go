package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Chat struct {
	TelegramChat  *tgbotapi.Chat
	Type          string
	Title         string
	Championships []*Championship
}
