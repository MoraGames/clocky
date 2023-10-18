package model

import "time"

type Championship struct {
	Title     string
	StartDate time.Time
	Duration  time.Duration
	Ranking   []Placement
}
