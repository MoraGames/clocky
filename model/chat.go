package model

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Chat struct {
		Id            int64
		TelegramChat  *tgbotapi.Chat
		Type          string
		Title         string
		Settings      *ChatSettings
		Championships []*Championship
	}

	ChatSettings struct {
		SummaryStats     bool
		SummaryStatsTime time.Duration
	}
)
