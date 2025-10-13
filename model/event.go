package model

import "time"

type (
	Event struct {
		Sets           []*Set
		Message        string
		Time           time.Time
		BasePoints     int
		Partecipations []*EventPartecipation
	}

	EventPartecipation struct {
		User         *User
		Championship *Championship
		Time         time.Time
	}
)
