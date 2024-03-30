package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	User struct {
		Id           int64
		TelegramUser *tgbotapi.User
		Nickname     string
		Stats        *UserStats
		Inventory    *UserInventory
	}

	UserStats struct {
		TotalPoints                      int
		TotalEventsWins                  int
		TotalEventsPartecipations        int
		TotalChampionshipsWins           int
		TotalChampionshipsPartecipations int
		MaxPointsInAChampionship         int
		MaxWinsInAChampionship           int
		MaxPartecipationsInAChampionship int
	}

	UserInventory struct {
		Gems    int
		Items   [5]*Effect
		Effects []*Effect
	}
)
