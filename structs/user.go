package structs

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
			if e != effect {
				newUserEffect = append(newUserEffect, e)
			}
		}
	}
	u.Effects = newUserEffect
}
