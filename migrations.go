package main

import (
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func manageMigrations() {
	fixBrokenActivityStreakBonusesEffectsMigration()
}

func fixBrokenActivityStreakBonusesEffectsMigration() {
	for userId := range Users {
		Users[userId].RemoveEffect(structs.NoNegative)
		Users[userId].RemoveEffect(structs.ConsistentParticipant1)
		Users[userId].RemoveEffect(structs.ConsistentParticipant2)
		if Users[userId].DailyActivityStreak >= 21 {
			Users[userId].AddEffect(structs.NoNegative)
		} else if Users[userId].DailyActivityStreak >= 14 {
			Users[userId].AddEffect(structs.ConsistentParticipant2)
		} else if Users[userId].DailyActivityStreak >= 7 {
			Users[userId].AddEffect(structs.ConsistentParticipant1)
		}
	}
}

func AddTelegramUserToExistingUser(msgFrom *tgbotapi.User) {
	if user, exist := Users[msgFrom.ID]; exist {
		user.TelegramUser = msgFrom
		Users[msgFrom.ID] = user
	}
}
