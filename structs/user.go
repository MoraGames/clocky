package structs

import "fmt"

type User struct {
	TelegramID                      int64
	UserName                        string
	TotalPoints                     int
	TotalEventPartecipations        int
	TotalEventWins                  int
	TotalChampionshipPartecipations int
	TotalChampionshipWins           int
	Effects                         []*Effect
}

func NewUser(telegramID int64, username string) *User {
	return &User{telegramID, username, 0, 0, 0, 0, 0, make([]*Effect, 0)}
}

func (u *User) AddEffects(effects ...*Effect) {
	u.Effects = append(u.Effects, effects...)
}

func (u *User) RemoveEffects(effects ...*Effect) {
	newUserEffect := make([]*Effect, 0)
	for _, e := range u.Effects {
		for _, effect := range effects {
			if e.Name != effect.Name {
				newUserEffect = append(newUserEffect, e)
			}
		}
	}
	u.Effects = newUserEffect
}

func (u *User) StringifyEffects() string {
	stringifiedEffects := ""
	for i, e := range u.Effects {
		if i != len(u.Effects)-1 {
			stringifiedEffects += fmt.Sprintf("%q, ", e.Name)
		} else {
			stringifiedEffects += fmt.Sprintf("%q", e.Name)
		}
	}
	return "[" + stringifiedEffects + "]"
}
