package model

import "time"

type Event struct {
	Message string
	Time    time.Time
	Points  int
	Bonus   *Bonus
}
