package model

import "time"

type (
	Championship struct {
		ID           int64
		Title        string
		StartDate    time.Time
		Settings     *ChampionshipSettings
		Partecipants []*User
		Ranking      []*ChampionshipPlacement
	}

	ChampionshipSettings struct {
		JoiningAvaibility               string
		JoiningParameter                string
		EndingType                      string
		EndingParameter                 string
		SetsEnabled                     []*Set
		SetsRotation                    string
		SetsRotationTime                time.Duration
		EffectsEnabled                  string
		EffectsEventMinimumAmount       int
		EffectsEventMaximumAmount       int
		EffectsInventoryAmount          int
		EffectsInventoryStackAmount     int
		EffectsInventoryShopEnabled     bool
		EffectsInventoryShopAmount      int
		EffectsInventoryShopRefreshTime time.Duration
	}

	ChampionshipPlacement struct {
		User                *User
		Position            int
		Points              int
		EventPartecipations int
		EventWins           int
	}
)
