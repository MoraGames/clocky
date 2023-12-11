package controller

import (
	"time"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
)

func (c *Controller) CreateEvent(message string) error {
	//Check if the user already exists
	if _, err := c.event.Get(message); err == nil {
		return errorType.ErrEventAlreadyExist{
			EventMessage: message,
			Message:      "cannot create event that already exists",
			Location:     "Controller.CreateEvent()",
		}
	} else if err.Error() != "cannot get event not found" {
		return err
	}

	time, err := time.Parse("15:04", message)
	if err != nil {
		return err
	}

	//Create the user
	event := &model.Event{
		Message: message,
		Time:    time,
		Points:  0,
		Effects: nil,
	}

	return c.event.Create(event)
}

func (c *Controller) GetEvent(message string) (*model.Event, error) {
	return c.event.Get(message)
}

func (c *Controller) GetAllEvents() []*model.Event {
	return c.event.GetAll()
}

func (c *Controller) DeleteEvent(message string) error {
	//Check if the chat already exists
	_, err := c.event.Get(message)
	if err != nil {
		return err
	}

	//Delete the chat
	return c.event.Delete(message)
}

func (c *Controller) IsEvent(message string) (bool, string, error) {
	//TODO: Implements IsEvent function
	return false, "", nil
}
