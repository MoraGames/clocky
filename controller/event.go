package controller

import (
	"time"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
)

func (c *Controller) CreateEvent(eventMessage, eventType string, eventTime time.Time) error {
	//Check if the user already exists
	if _, err := c.event.Get(eventMessage); err == nil {
		return errorType.ErrEventAlreadyExist{
			EventMessage: eventMessage,
			Message:      "cannot create event that already exists",
			Location:     "Controller.CreateEvent()",
		}
	} else if err.Error() != "cannot get event not found" {
		return err
	}

	if eventType == "time" && (eventMessage != eventTime.Format("15:04")) {
		return errorType.ErrEventNotValid{
			EventMessage: eventMessage,
			Message:      "cannot create event of type \"time\" with different message and time.Format(\"15:04\") properties.",
			Location:     "Controller.CreateEvent()",
		}
	} else if eventType == "second" && (eventMessage != eventTime.Format("15:04:05")) {
		return errorType.ErrEventNotValid{
			EventMessage: eventMessage,
			Message:      "cannot create event of type \"second\" with different message and time.Format(\"15:04:05\") properties.",
			Location:     "Controller.CreateEvent()",
		}
	}

	//Create the user
	event := &model.Event{
		Message:   eventMessage,
		Type:      eventType,
		Time:      eventTime,
		Enabled:   false,
		Instances: nil,
	}

	return c.event.Create(event)
}

func (c *Controller) GetEvent(message string) (*model.Event, error) {
	return c.event.Get(message)
}

func (c *Controller) GetAllEvents() []*model.Event {
	return c.event.GetAll()
}

func (c *Controller) UpdateEvent(message string, event *model.Event) error {
	//Check if the event already exists
	_, err := c.event.Get(message)
	if err != nil {
		return err
	}

	//Update the event
	return c.event.Update(message, event)
}

func (c *Controller) DeleteEvent(message string) error {
	//Check if the event already exists
	_, err := c.event.Get(message)
	if err != nil {
		return err
	}

	//Delete the event
	return c.event.Delete(message)
}

func (c *Controller) IsEvent(message string) (bool, *model.Event, error) {
	//Check if the event exists
	event, err := c.event.Get(message)
	if err != nil {
		if err.Error() == "cannot get event not found" {
			return false, nil, nil
		}
		return false, nil, err
	}

	//Return the data
	return true, event, nil
}
