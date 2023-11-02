package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.EventRepoer = new(EventRepo)

// mock.UserRepo
type EventRepo struct {
	events map[string]*model.Event
}

// Return a new UserRepo
func NewEventRepo() *EventRepo {
	return &EventRepo{
		events: make(map[string]*model.Event),
	}
}

func (er *EventRepo) Create(event *model.Event) error {
	if _, ok := er.events[event.Message]; ok {
		return errorType.ErrEventAlreadyExist{
			EventMessage: event.Message,
			Message:      "cannot create event that already exists",
			Location:     "EventRepo.Create()",
		}
	}

	er.events[event.Message] = event
	return nil
}

func (er *EventRepo) Get(message string) (*model.Event, error) {
	record, ok := er.events[message]
	if !ok {
		return nil, errorType.ErrEventNotFound{
			EventMessage: message,
			Message:      "cannot get event not found",
			Location:     "EventRepo.Get()",
		}
	}
	return record, nil
}

func (er *EventRepo) GetAll() []*model.Event {
	events := make([]*model.Event, 0, len(er.events))
	for _, event := range er.events {
		events = append(events, event)
	}
	return events
}

func (er *EventRepo) Update(message string, event *model.Event) error {
	if _, ok := er.events[message]; !ok {
		return errorType.ErrEventNotFound{
			EventMessage: message,
			Message:      "cannot update event not found",
			Location:     "EventRepo.Update()",
		}
	}
	if message != event.Message {
		return errorType.ErrEventNotValid{
			EventMessage: message,
			Message:      "cannot update event when message mismatch",
			Location:     "EventRepo.Update()",
		}
	}

	er.events[message] = event
	return nil
}

func (er *EventRepo) Delete(message string) error {
	if _, ok := er.events[message]; !ok {
		return errorType.ErrEventNotFound{
			EventMessage: message,
			Message:      "cannot delete event not found",
			Location:     "EventRepo.Delete()",
		}
	}
	delete(er.events, message)
	return nil
}
