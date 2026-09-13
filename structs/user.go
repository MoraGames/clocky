package structs

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	TelegramUser                    *tgbotapi.User
	TelegramID                      int64
	UserName                        string
	TotalPoints                     int
	TotalEventPartecipations        int
	TotalEventWins                  int
	DailyPoints                     int
	DailyEventPartecipations        int
	DailyEventWins                  int
	ChampionshipPoints              int
	ChampionshipEventPartecipations int
	ChampionshipEventWins           int
	TotalChampionshipPartecipations int
	TotalChampionshipWins           int
	DailyPartecipationStreak        int
	DailyActivityStreak             int
	ChampionshipParticipationStreak int
	ChampionshipWinStreak           int
	ChampionshipAbsenceStreak       int
	OverheatingPoints               int
	Overheating                     bool
	LastEventSequence               int64
	Effects                         []*Effect
	FirstParticipation              time.Time
}

type UserMinimal struct {
	TelegramID int64
	UserName   string
}

// DisplayName returns the @username when the user has one, falling back to
// their first and last name otherwise. Users without a @username (only a
// display name set) would otherwise be shown and stored with an empty name.
func DisplayName(telegramUser *tgbotapi.User) string {
	if telegramUser.UserName != "" {
		return telegramUser.UserName
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", telegramUser.FirstName, telegramUser.LastName))
}

func NewUser(telegramUser *tgbotapi.User) *User {
	return &User{TelegramUser: telegramUser, TelegramID: telegramUser.ID, UserName: DisplayName(telegramUser), Effects: make([]*Effect, 0), FirstParticipation: time.Now()}
}

func (u *User) RegisterEventParticipation(eventSequence int64) bool {
	if eventSequence <= u.LastEventSequence {
		return true
	}

	for missed := u.LastEventSequence + 1; missed < eventSequence; missed++ {
		u.ChampionshipParticipationStreak = 0
		u.ChampionshipWinStreak = 0
		u.ChampionshipAbsenceStreak++
		u.OverheatingPoints += max(-2-(u.ChampionshipAbsenceStreak/12), -4)
	}

	u.ChampionshipParticipationStreak++
	u.ChampionshipWinStreak = 0
	u.ChampionshipAbsenceStreak = 0
	u.LastEventSequence = eventSequence
	u.updateOverheating()
	return true
}

func (u *User) RegisterEventWin() {
	u.ChampionshipWinStreak++
	u.OverheatingPoints += 3 + (u.ChampionshipWinStreak / 6)
	u.updateOverheating()
}

func (u *User) ResetChampionshipOverheating() {
	u.ChampionshipParticipationStreak = 0
	u.ChampionshipWinStreak = 0
	u.ChampionshipAbsenceStreak = 0
	u.OverheatingPoints = 0
	u.Overheating = false
}

func (u *User) IsOverheating() bool {
	u.updateOverheating()
	return u.Overheating
}

func (u *User) updateOverheating() {
	if u.Overheating {
		if u.OverheatingPoints < 40 {
			u.Overheating = false
		}
		return
	}
	if u.OverheatingPoints > 90 {
		u.Overheating = true
	}
}

func (u *User) Minimize() *UserMinimal {
	return &UserMinimal{u.TelegramID, u.UserName}
}

func (u *User) AddEffect(effectToAdd *Effect) {
	u.Effects = append(u.Effects, effectToAdd)
}

func (u *User) RemoveEffect(effectToRemove *Effect) {
	newUserEffects := make([]*Effect, 0)
	for _, userEffect := range u.Effects {
		if userEffect.Name != effectToRemove.Name {
			newUserEffects = append(newUserEffects, userEffect)
		}
	}
	u.Effects = newUserEffects
}

func (u *User) StringifyEffects(brackets bool) string {
	stringifiedEffects := ""
	for i, e := range u.Effects {
		if i != len(u.Effects)-1 {
			stringifiedEffects += fmt.Sprintf("%q, ", e.Name)
		} else {
			stringifiedEffects += fmt.Sprintf("%q", e.Name)
		}
	}
	if brackets {
		return "[" + stringifiedEffects + "]"
	}
	return stringifiedEffects
}
