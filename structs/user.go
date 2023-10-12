package structs

type User struct {
	UserName                        string
	TotalPoints                     int
	TotalEventPartecipations        int
	TotalEventWins                  int
	TotalChampionshipPartecipations int
	TotalChampionshipWins           int
}

func NewUser(username string) *User {
	return &User{username, 0, 0, 0, 0, 0}
}
