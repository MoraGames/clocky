package model

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Chat struct {
		TelegramChat        *tgbotapi.Chat
		Settings            *ChatSettings
		Championships       []*Championship
		CurrentChampionship int
	}

	ChatSettings struct {
		SummaryEnabled   []ChatSummarySubject
		SummaryFrequency time.Duration
		SummaryLocation  string // "group"|"private"
	}
	ChatSummarySubject string
)

const (
	SummaryStats  ChatSummarySubject = "stats"
	SummaryHints  ChatSummarySubject = "hints"
	SummaryTop3   ChatSummarySubject = "top3"
	SummaryTop5   ChatSummarySubject = "top5"
	SummaryTop10  ChatSummarySubject = "top10"
	SummaryActive ChatSummarySubject = "active"
)
