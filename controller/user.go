package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Controller) CreateUser(userID int64, telegramUser *tgbotapi.User) error {
	//Check if the user already exists
	if _, err := c.user.Get(userID); err == nil {
		return errorType.ErrUserAlreadyExist{
			UserID:   userID,
			Message:  "cannot create user that already exists",
			Location: "UserController.CreateUser()",
		}
	} else if err.Error() != "cannot get user not found" {
		return err
	}

	//Create the user
	user := &model.User{
		TelegramUser: telegramUser,
		UserStats: &model.UserStats{
			TotalPoints:                      0,
			MaxChampionshipPoints:            0,
			MaxChampionshipWins:              0,
			TotalEventsPartecipations:        0,
			TotalChampionshipsPartecipations: 0,
			TotalEventsWins:                  0,
			TotalChampionshipsWins:           0,
		},
	}

	return c.user.Create(user)
}

func (c *Controller) GetUser(userID int64) (*model.User, error) {
	return c.user.Get(userID)
}

func (c *Controller) GetAllUsers() ([]*model.User, error) {
	return c.user.GetAll()
}

func (c *Controller) GetUserStats(userID int64) (*model.UserStats, error) {
	//Check if the user already exists
	user, err := c.user.Get(userID)
	if err != nil {
		return nil, err
	}

	return user.UserStats, nil
}

func (c *Controller) GetUserTotalPoints(userID int64) (int, error) {
	//Check if the user already exists
	user, err := c.user.Get(userID)
	if err != nil {
		return 0, err
	}

	return user.UserStats.TotalPoints, nil
}

func (c *Controller) SetUserTotalPoints(userID int64, userPoints int) error {
	//Check if the user already exists
	user, err := c.user.Get(userID)
	if err != nil {
		return err
	}

	//Check if the userPoints are valid
	if userPoints < 0 {
		return errorType.ErrNegativeNumber{
			Number:   userPoints,
			Message:  "userPoints cannot be negative",
			Location: "UserController.SetUserPoints()",
		}
	}

	//Update the user
	user.UserStats.TotalPoints = userPoints

	return c.user.Update(userID, user)
}

func (c *Controller) UpdateUser(userID int64, user *model.User) error {
	//Check if the user already exists
	_, err := c.user.Get(userID)
	if err != nil {
		return err
	}

	//Update the user
	return c.user.Update(userID, user)
}

func (c *Controller) ResetUser(userID int64) error {
	//Check if the user already exists
	user, err := c.user.Get(userID)
	if err != nil {
		return err
	}

	//Reset the user
	user = &model.User{
		TelegramUser: user.TelegramUser,
		UserStats: &model.UserStats{
			TotalPoints:                      0,
			MaxChampionshipPoints:            0,
			MaxChampionshipWins:              0,
			TotalEventsPartecipations:        0,
			TotalChampionshipsPartecipations: 0,
			TotalEventsWins:                  0,
			TotalChampionshipsWins:           0,
		},
	}

	return c.user.Update(userID, user)
}

func (c *Controller) DeleteUser(userID int64) error {
	//Check if the user already exists
	_, err := c.user.Get(userID)
	if err != nil {
		return err
	}

	//Delete the user
	return c.user.Delete(userID)
}
