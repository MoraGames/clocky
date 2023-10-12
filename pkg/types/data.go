package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Data struct {
	Bot *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
}