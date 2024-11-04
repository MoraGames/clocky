package model

import (
	"time"
)

type (
	Event struct {
		ID        int64
		Sets      []*Set
		Message   string
		Type      string
		Time      time.Time
		Enabled   bool
		Instances []*EventInstance
	}
	EventInstance struct {
		ID             int64
		Event          *Event
		BasePoints     int
		Points         int
		Effects        []*Effect
		Activated      bool
		ActivatedBy    *User
		ActivatedAt    time.Time
		ArrivedAt      time.Time
		Partecipations []*EventPartecipation
	}
	EventPartecipation struct {
		ID            int64
		EventInstance *EventInstance
		User          *User
		Chat          *Chat
		Time          time.Time
	}
)
