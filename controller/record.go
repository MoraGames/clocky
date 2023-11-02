package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
)

func (c *Controller) CreateRecord(title string) error {
	//Check if the user already exists
	if _, err := c.record.Get(title); err == nil {
		return errorType.ErrRecordAlreadyExist{
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

	return c.record.Create(record)
}

func (c *Controller) GetRecord(title string) (*model.Record, error) {
	return c.record.Get(title)
}

func (c *Controller) GetAllRecords() []*model.Record {
	return c.record.GetAll()
}

func (c *Controller) DeleteRecord(title string) error {
	//Check if the chat already exists
	_, err := c.record.Get(title)
	if err != nil {
		return err
	}

	//Delete the chat
	return c.record.Delete(title)
}
