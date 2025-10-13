package model

import "time"

type (
	Championship struct {
		ID        int64
		Title     string
		StartDate time.Time
		Settings  *ChampionshipSettings
		Ranking   []*ChampionshipPlacement // Sorted by Points > Wins > Partecipations > User.ID
	}

	ChampionshipSettings struct {
		JoiningAvaibility string // "public"|"private"
		JoiningWindow     string // "anytime"|"scheduled"
		JoiningParameter  string // ^ ""|duration
		EndingType        string // "endless"|"points"|"gaps"|"duration"|"events"
		EndingParameter   string // ^ ""|points|gap|duration|events
	}

	ChampionshipPlacement struct {
		User                *User
		Points              int
		EventPartecipations int
		EventWins           int
	}
)
