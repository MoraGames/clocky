package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	User struct {
		ID           int64
		TelegramUser *tgbotapi.User
		Nickname     string
		Settings     *UserSettings
		Stats        *UserStats
		MaxStats     *UserMaxStats
		Inventory    *UserInventory
		Tracking     *UserTracking
	}

	UserSettings struct {
		Language                 string
		TotalStatsPrivacy        string
		ChampionshipStatsPrivacy string
		DailyStatsPrivacy        string
		RotationStatsPrivacy     string
	}

	UserStats struct {
		//Event's Points/Wins/Partecipations in the lifetime of the day
		DailyPoints               int
		DailyEventsWins           int
		DailyEventsPartecipations int
		//Event's Points/Wins/Partecipations in the lifetime of the rotation
		RotationPoints               int
		RotationEventsWins           int
		RotationEventsPartecipations int
		//Event's Points/Wins/Partecipations in the lifetime of the championship
		ChampionshipPoints               int
		ChampionshipEventsWins           int
		ChampionshipEventsPartecipations int
		//Event's Points/Wins/Partecipations in the lifetime of the user
		TotalPoints               int
		TotalEventsWins           int
		TotalEventsPartecipations int
		//Championship's Wins/Partecipations in the lifetime of the user
		TotalChampionshipsWins           int
		TotalChampionshipsPartecipations int
	}
	UserMaxStats struct {
		//Max Event's Points/Wins/Partecipations in the lifetime of the day
		MaxDailyPoints               int
		MaxDailyEventsWins           int
		MaxDailyEventsPartecipations int
		//Event's Points/Wins/Partecipations in the lifetime of the rotation
		MaRotationPoints               int
		MaRotationEventsWins           int
		MaRotationEventsPartecipations int
		//Event's Points/Wins/Partecipations in the lifetime of the championship
		MaxChampionshipPoints               int
		MaxChampionshipEventsWins           int
		MaxChampionshipEventsPartecipations int
	}

	UserInventory struct {
		Gems    int
		Items   [5]*Effect
		Effects []*Effect
	}

	UserTracking struct {
		//TODO: Implements anti-cheat tracking system
	}
)
