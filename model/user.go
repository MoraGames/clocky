package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	User struct {
		TelegramUser *tgbotapi.User
		UserStats    *UserStats
	}

	UserStats struct {
		TotalPoints                      int
		MaxChampionshipPoints            int
		TotalEventPartecipations         int
		TotalEventWins                   int
		TotalChampionshipsPartecipations int
		TotalChampionshipsWins           int
	}
)
