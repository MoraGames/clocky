package model

type (
	//TODO: Rework this shit
	Effect struct {
		ID         int64
		Name       string
		Type       string
		Parameters map[string]int
	}
)
