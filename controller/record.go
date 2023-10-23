package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
	"github.com/sirupsen/logrus"
)

type RecordController struct {
	repo repo.RecordRepoer
	log  *logrus.Logger
}

func NewRecordController(repoer repo.ChatRepoer, logger *logrus.Logger) *ChatController {
	return &ChatController{
		repo: repoer,
		log:  logger,
	}
}

func (rc *RecordController) CreateRecord(title string) error {
	//Check if the user already exists
	if _, err := rc.repo.Get(title); err == nil {
		return errorType.ErrRecordAlreadyExists{
			RecordTitle: title,
			Message:     "cannot create record that already exists",
			Location:    "RecordController.CreateRecord()",
		}
	} else if err.Error() != "cannot get record not found" {
		return err
	}

	//Create the user
	record := &model.Record{
		Title:        title,
		Value:        0,
		User:         nil,
		Championship: nil,
	}

	return rc.repo.Create(record)
}

func (rc *RecordController) GetRecord(title string) (*model.Record, error) {
	return rc.repo.Get(title)
}

func (rc *RecordController) GetAllRecords() []*model.Record {
	return rc.repo.GetAll()
}

func (rc *RecordController) DeleteRecord(title string) error {
	//Check if the chat already exists
	_, err := rc.repo.Get(title)
	if err != nil {
		return err
	}

	//Delete the chat
	return rc.repo.Delete(title)
}
