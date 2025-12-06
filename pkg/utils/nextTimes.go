package utils

import "time"

func NextInWeekdayAtTime(t time.Time, weekday time.Weekday, hour, min, sec int) time.Time {
	// Build today's date at the target time
	target := time.Date(
		t.Year(), t.Month(), t.Day(),
		hour, min, sec, 0,
		t.Location(),
	)
	return target
}
