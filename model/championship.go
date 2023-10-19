package model

import "time"

type Championship struct {
	ID        int64
	Title     string
	StartDate time.Time
	Duration  time.Duration
	Ranking   []Placement
}
