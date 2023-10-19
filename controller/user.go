package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	repo repo.UserRepoer
	log  *logrus.Logger
}

func NewUserController(repoer repo.UserRepoer, logger *logrus.Logger) *UserController {
	return &UserController{
		repo: repoer,
		log:  logger,
	}
}

func (uc *UserController) CreateUser(userID int64, telegramUser *tgbotapi.User) error {
	//Check if the user already exists
	if _, err := uc.repo.Get(userID); err == nil {
		return errorType.ErrUserAlreadyExists{
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
			TotalEventPartecipations:         0,
			TotalEventWins:                   0,
			TotalChampionshipsPartecipations: 0,
			TotalChampionshipsWins:           0,
		},
	}

	return uc.repo.Create(user)
}

func (uc *UserController) GetUser(userID int64) (*model.User, error) {
	return uc.repo.Get(userID)
}

func (uc *UserController) GetAllUsers() []*model.User {
	return uc.repo.GetAll()
}

func (uc *UserController) GetUserStats(userID int64) (*model.UserStats, error) {
	//Check if the user already exists
	user, err := uc.repo.Get(userID)
	if err != nil {
		return nil, err
	}

	return user.UserStats, nil
}

func (uc *UserController) GetUserTotalPoints(userID int64) (int, error) {
	//Check if the user already exists
	user, err := uc.repo.Get(userID)
	if err != nil {
		return 0, err
	}

	return user.UserStats.TotalPoints, nil
}

func (uc *UserController) SetUserTotalPoints(userID int64, userPoints int) error {
	//Check if the user already exists
	user, err := uc.repo.Get(userID)
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

	return uc.repo.Update(userID, user)
}

func (uc *UserController) ResetUser(userID int64) error {
	//Check if the user already exists
	user, err := uc.repo.Get(userID)
	if err != nil {
		return err
	}

	//Reset the user
	user = &model.User{
		TelegramUser: user.TelegramUser,
		UserStats: &model.UserStats{
			TotalPoints:                      0,
			MaxChampionshipPoints:            0,
			TotalEventPartecipations:         0,
			TotalEventWins:                   0,
			TotalChampionshipsPartecipations: 0,
			TotalChampionshipsWins:           0,
		},
	}

	return nil
}

func (uc *UserController) DeleteUser(userID int64) error {
	//Check if the user already exists
	_, err := uc.repo.Get(userID)
	if err != nil {
		return err
	}

	//Delete the user
	return uc.repo.Delete(userID)
}
