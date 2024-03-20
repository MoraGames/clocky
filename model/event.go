package model

import (
	"time"
)

type (
	Event struct {
		ID      int64
		Message string
		Type    string
		Time    time.Time
		Enabled bool
	}
	EventInstance struct {
		ID          int64
		Event       *Event
		BasePoints  int
		Points      int
		Effects     []*Effect
		Activated   bool
		ActivatedBy *User
		ActivatedAt time.Time
		ArrivedAt   time.Time
	}
	EventPartecipation struct {
		ID            int64
		EventInstance *EventInstance
		User          *User
		Chat          *Chat
		Time          time.Time
	}
)
