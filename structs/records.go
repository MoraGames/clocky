package structs

/*
type RecordType string

// type RecordValue interface {
// 	GetUser() *User
// }

type AbsoluteRecord struct {
	User   *User
	Record int
}

// func (ar AbsoluteRecord) GetUser() *User {
// 	return ar.User
// }

type ChampionshipRecord struct {
	User         *User
	Championship *Championship
	Record       int
}

// func (cr ChampionshipRecord) GetUser() *User {
// 	return cr.User
// }

type RecordsMap[RecordValue any] map[RecordType]RecordValue

func main() {
	var Records = RecordsMap[interface{}]{
		"MostPoints":                             AbsoluteRecord{nil, 0},
		"MostPointsInAChampionship":              ChampionshipRecord{nil, nil, 0},
		"MostEventPartecipations":                AbsoluteRecord{nil, 0},
		"MostEventPartecipationsInAChampionship": ChampionshipRecord{nil, nil, 0},
		"MostEventWins":                          AbsoluteRecord{nil, 0},
		"MostEventWinsInAChampionship":           ChampionshipRecord{nil, nil, 0},
	}

	//OK  -->  fmt.Println(Records)
	//OK  -->  fmt.Println(Records["MostEventWinsInAChampionship"])
	//NO (RecordValue non definisce .Championship)  -->  fmt.Println(Records["MostEventWinsInAChampionship"].Championship)
}
*/
