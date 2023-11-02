package model

type Set struct {
	ID     int64
	Type   string
	Rule   string
	Events []*Event
}
