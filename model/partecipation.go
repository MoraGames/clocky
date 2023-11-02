package model

import "time"

type Partecipation struct {
	ID    int64
	User  *User
	Event *Event
	Chat  *Chat
	Time  time.Time
}
