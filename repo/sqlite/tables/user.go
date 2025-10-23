package tables

import (
	"github.com/MoraGames/clockyuwu/model"
	"gorm.io/gorm"
)

type UserTable struct {
	gorm.Model
	TelegramID       int64
	TelegramIsBot    bool
	TelegramUsername string
	Nickname         string
	Settings         UserSettingsTable
	Stats            UserStatsTable
	MaxStats         UserMaxStatsTable
}

func (UserTable) FromModel(user *model.User) *UserTable {
	return &UserTable{
		TelegramID: user.TelegramID,
		Nickname:   user.Nickname,
	}
}

func (UserTable) ToModel() *model.User {
	return &model.User{}
}

type UserSettingsTable struct {
	gorm.Model
	UserID uint
	model.UserSettings
}

func (UserSettingsTable) FromModel(settings *model.UserSettings) UserSettingsTable {
	return UserSettingsTable{
		UserSettings: *settings,
	}
}

func (UserSettingsTable) ToModel() *model.UserSettings {
	return &model.UserSettings{}
}

type UserStatsTable struct {
	gorm.Model
	UserID uint
	model.UserStats
}

type UserMaxStatsTable struct {
	gorm.Model
	UserID uint
	model.UserMaxStats
}

func (UserStatsTable) ToModel() *model.UserStats {
	return &model.UserStats{}
}

func (UserMaxStatsTable) ToModel() *model.UserMaxStats {
	return &model.UserMaxStats{}
}
