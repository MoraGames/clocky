package model

import (
	"time"
)

type (
	Event struct {
		Message   string
		Type      string
		Time      time.Time
		Enabled   bool
		Instances []*EventIstance
	}
	EventIstance struct {
		Points         int
		Effects        []*Effect
		Activated      bool
		ActivatedBy    *User
		ActivatedAt    time.Time
		ArrivedAt      time.Time
		Partecipations []EventPartecipation
	}
	EventPartecipation struct {
		User *User
		Time time.Time
	}
)
