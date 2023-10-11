package events

import "time"

type EventValue struct {
	Points      int
	Activated   bool
	ActivatedBy string
	ActivatedAt time.Time
}
