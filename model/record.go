package model

type Record struct {
	Title        string
	Value        int
	User         *User
	Championship *Championship
}

/*
var Records = RecordsMap[interface{}]{
		"MostPoints":                             AbsoluteRecord{nil, 0},
		"MostPointsInAChampionship":              ChampionshipRecord{nil, nil, 0},
		"MostEventPartecipations":                AbsoluteRecord{nil, 0},
		"MostEventPartecipationsInAChampionship": ChampionshipRecord{nil, nil, 0},
		"MostEventWins":                          AbsoluteRecord{nil, 0},
		"MostEventWinsInAChampionship":           ChampionshipRecord{nil, nil, 0},
	}
*/
