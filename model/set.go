package model

type Set struct {
	Id     int64
	Name   string
	Type   string
	Rule   string
	Events []*Event
}
