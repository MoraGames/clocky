package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	User struct {
		ID           int64
		TelegramUser *tgbotapi.User
		Nickname     string
		UserStats    *UserStats
		Effects      []*Effect
	}

	UserStats struct {
		TotalPoints                      int
		MaxChampionshipPoints            int
		MaxChampionshipWins              int
		TotalEventsPartecipations        int
		TotalChampionshipsPartecipations int
		TotalEventsWins                  int
		TotalChampionshipsWins           int
	}
)
