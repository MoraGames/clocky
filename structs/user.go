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
	Overheating                     bool
	OverheatingParticipationStreak  int
	OverheatingMissedEventStreak    int
	OverheatingLastEventSequence    int64
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

const (
	overheatingParticipationLimit = 15
	overheatingRecoveryEvents     = 5
)

func (u *User) RegisterEventParticipation(eventSequence int64) bool {
	if eventSequence <= u.OverheatingLastEventSequence {
		return !u.Overheating
	}

	missedEvents := int(eventSequence - u.OverheatingLastEventSequence - 1)
	if u.Overheating {
		if missedEvents >= overheatingRecoveryEvents {
			u.Overheating = false
			u.OverheatingParticipationStreak = 0
			u.OverheatingMissedEventStreak = 0
		} else {
			u.OverheatingMissedEventStreak += missedEvents
			u.OverheatingLastEventSequence = eventSequence
			return false
		}
	}

	if missedEvents == 0 && u.OverheatingLastEventSequence != 0 {
		u.OverheatingParticipationStreak++
	} else {
		u.OverheatingParticipationStreak = 1
	}
	u.OverheatingLastEventSequence = eventSequence
	if u.OverheatingParticipationStreak > overheatingParticipationLimit {
		u.Overheating = true
	}
	return true
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
