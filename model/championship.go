package model

import "time"

type (
	Championship struct {
		ID        int64
		Title     string
		StartDate time.Time
		Duration  time.Duration
		Ranking   []*ChampionshipPlacement
	}

	ChampionshipPlacement struct {
		User                *User
		Position            int
		Points              int
		EventPartecipations int
		EventWins           int
	}
)
