package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	User struct {
		TelegramUser *tgbotapi.User
		UserStats    *UserStats
		Effects      []*Effect
	}

	UserStats struct {
		TotalPoints                      int
		MaxChampionshipPoints            int
		MaxChampionshipWins              int
		TotalEventPartecipations         int
		TotalChampionshipsPartecipations int
		TotalEventWins                   int
		TotalChampionshipsWins           int
	}
)
