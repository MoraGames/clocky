package mock

import (
	"context"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// // UserRepo Error
// type ErrUserRepo struct {
// 	UserId   int64
// 	Message  string
// 	Location string
// }

// func (err ErrUserRepo) Error() string {
// 	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.UserId)
// }

// Check if the repo implements the interface
var _ repo.UserRepoer = new(UserRepo)

// UserRepo is a mock implementation
type UserRepo struct {
	dbLocation string
}

func NewUserRepo(dbLocation string) *UserRepo {
	return &UserRepo{
		dbLocation: dbLocation,
	}
}

func (ur *UserRepo) Open() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(ur.dbLocation), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ur *UserRepo) Create(ctx context.Context, user *model.User) (uint, error) {
	db, err := ur.Open()
	if err != nil {
		return 0, err
	}

	_, err = gorm.G[UserTable](db).Select("TelegramUserID").Where("TelegramUserID = ?", user.TelegramUser.ID).First(ctx)
	if err == nil || err != gorm.ErrRecordNotFound {
		return 0, errorType.ErrUserAlreadyExist{
			UserID:   user.TelegramUser.ID,
			Message:  "user already exist",
			Location: "UserRepo.Create()",
		}
	}

	var newUser = UserTable{
		TelegramUserID: user.TelegramUser.ID,
		Nickname:       user.Nickname,
		Settings: UserSettingsTable{
			UserSettings: model.UserSettings{
				DailyStatsPrivacy:        user.Settings.DailyStatsPrivacy,
				RotationStatsPrivacy:     user.Settings.RotationStatsPrivacy,
				ChampionshipStatsPrivacy: user.Settings.ChampionshipStatsPrivacy,
				TotalStatsPrivacy:        user.Settings.TotalStatsPrivacy,
			},
		},
		Stats: UserStatsTable{
			UserStats: model.UserStats{
				DailyPoints:                      user.Stats.DailyPoints,
				DailyEventsWins:                  user.Stats.DailyEventsWins,
				DailyEventsPartecipations:        user.Stats.DailyEventsPartecipations,
				RotationPoints:                   user.Stats.RotationPoints,
				RotationEventsWins:               user.Stats.RotationEventsWins,
				RotationEventsPartecipations:     user.Stats.RotationEventsPartecipations,
				ChampionshipPoints:               user.Stats.ChampionshipPoints,
				ChampionshipEventsWins:           user.Stats.ChampionshipEventsWins,
				ChampionshipEventsPartecipations: user.Stats.ChampionshipEventsPartecipations,
				TotalPoints:                      user.Stats.TotalPoints,
				TotalEventsWins:                  user.Stats.TotalEventsWins,
				TotalEventsPartecipations:        user.Stats.TotalEventsPartecipations,
				TotalChampionshipsWins:           user.Stats.TotalChampionshipsWins,
				TotalChampionshipsPartecipations: user.Stats.TotalChampionshipsPartecipations,
			},
		},
		MaxStats: UserMaxStatsTable{
			UserMaxStats: model.UserMaxStats{
				MaxDailyPoints:                      user.MaxStats.MaxDailyPoints,
				MaxDailyEventsWins:                  user.MaxStats.MaxDailyEventsWins,
				MaxDailyEventsPartecipations:        user.MaxStats.MaxDailyEventsPartecipations,
				MaxRotationPoints:                   user.MaxStats.MaxRotationPoints,
				MaxRotationEventsWins:               user.MaxStats.MaxRotationEventsWins,
				MaxRotationEventsPartecipations:     user.MaxStats.MaxRotationEventsPartecipations,
				MaxChampionshipPoints:               user.MaxStats.MaxChampionshipPoints,
				MaxChampionshipEventsWins:           user.MaxStats.MaxChampionshipEventsWins,
				MaxChampionshipEventsPartecipations: user.MaxStats.MaxChampionshipEventsPartecipations,
			},
		},
	}
	err = gorm.G[UserTable](db).Create(ctx, &newUser)
	return newUser.ID, nil
}

func (ur *UserRepo) Get(id int64) (*model.User, error) {
	user, ok := ur.users[id]
	if !ok {
		return nil, ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Get()",
		}
	}
	return user, nil
}

func (ur *UserRepo) GetAll() []*model.User {
	users := make([]*model.User, 0, len(ur.users))
	for _, user := range ur.users {
		users = append(users, user)
	}
	return users
}

func (ur *UserRepo) Update(id int64, user *model.User) error {
	_, ok := ur.users[id]
	if !ok {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Update()",
		}
	}
	if id != user.ID {
		return ErrUserRepo{
			UserId:   id,
			Message:  "users id mismatch",
			Location: "UserRepo.Update()",
		}
	}

	if ur.nicknameAlreadyUsed(user) {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user nickname already used",
			Location: "UserRepo.Update()",
		}
	}

	ur.users[id] = user
	return nil
}

func (ur *UserRepo) Delete(id int64) error {
	_, ok := ur.users[id]
	if !ok {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Delete()",
		}
	}
	delete(ur.users, id)
	return nil
}

func (ur *UserRepo) nicknameAlreadyUsed(user *model.User) bool {
	for _, u := range ur.users {
		if u.Nickname == user.Nickname {
			return true
		}
	}
	return false
}
