package utils

import (
	"fmt"
	"sort"

	"github.com/MoraGames/clockyuwu/structs"
)

// Rank is the type used for manage the leaderboard
type Rank struct {
	UserTelegramID int64
	Username       string
	Points         int
	Partecipations int
}

func GetRanking(Users map[int64]*structs.User) []Rank {
	// Generate the ranking
	ranking := make([]Rank, 0)
	for _, u := range Users {
		if u != nil {
			ranking = append(ranking, Rank{u.TelegramID, u.UserName, u.ChampionshipPoints, u.TotalEventPartecipations})
		}
	}

	// Sort the ranking by points (and partecipations if points are equal)
	sort.Slice(
		ranking,
		func(i, j int) bool {
			if ranking[i].Points == ranking[j].Points {
				return ranking[i].Partecipations < ranking[j].Partecipations
			}
			return ranking[i].Points > ranking[j].Points
		},
	)

	fmt.Printf(">>>>>>> Users: %+v\n", Users)
	fmt.Printf(">>>>>>> Generated ranking: %+v\n", ranking)

	return ranking
}
