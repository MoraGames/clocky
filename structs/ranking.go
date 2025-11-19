package structs

import (
	"fmt"
	"sort"
)

// RankScope is the type used for define the ranking scope (daily, championship, total) during the /ranking command execution
type RankScope string

const (
	RankScopeDay          RankScope = "day"
	RankScopeChampionship RankScope = "championship"
	RankScopeTotal        RankScope = "total"

	// DefaultRankScope is the default ranking scope used when no scope is provided
	DefaultRankScope = RankScopeChampionship
)

// Rank is the type used for manage the leaderboard
type Rank struct {
	UserTelegramID int64
	Username       string
	Points         int
	Partecipations int
}

func GetRanking(Users map[int64]*User, scope RankScope) []Rank {
	// Generate the ranking
	ranking := make([]Rank, 0)
	for _, u := range Users {
		if u != nil {
			switch scope {
			case RankScopeDay:
				ranking = append(ranking, Rank{u.TelegramID, u.UserName, u.DailyPoints, u.DailyEventPartecipations})
			case RankScopeChampionship:
				ranking = append(ranking, Rank{u.TelegramID, u.UserName, u.ChampionshipPoints, u.ChampionshipEventPartecipations})
			case RankScopeTotal:
				ranking = append(ranking, Rank{u.TelegramID, u.UserName, u.TotalPoints, u.TotalEventPartecipations})
			default:
				ranking = append(ranking, Rank{u.TelegramID, u.UserName, u.ChampionshipPoints, u.ChampionshipEventPartecipations})
			}
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
