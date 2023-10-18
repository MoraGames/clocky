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
	if _, err := uc.repo.Get(userID); err != nil {
		return err
	}

	//Create the user
	user := &model.User{
		TelegramUser:        telegramUser,
		Points:              0,
		EventPartecipations: 0,
		EventWins:           0,
	}

	return uc.repo.Create(user)
}

func (uc *UserController) GetUser(userID int64) (*model.User, error) {
	return uc.repo.Get(userID)
}

func (uc *UserController) GetAllUsers() []*model.User {
	return uc.repo.GetAll()
}

func (uc *UserController) GetUserPoints(userID int64) (int, error) {
	//Check if the user already exists
	user, err := uc.repo.Get(userID)
	if err != nil {
		return 0, err
	}

	return user.Points, nil
}

func (uc *UserController) SetUserPoints(userID int64, userPoints int) error {
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
	user.Points = userPoints

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
		TelegramUser:        user.TelegramUser,
		Points:              0,
		EventPartecipations: 0,
		EventWins:           0,
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
