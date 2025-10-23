package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// EventRepo Error
type ErrEventRepo struct {
	EventId  int64
	Message  string
	Location string
}

func (err ErrEventRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.EventId)
}

// Check if the repo implements the interface
var _ repo.EventRepoer = new(EventRepo)

// EventRepo is a mock implementation
type EventRepo struct {
	filePath string
}

func NewEventRepo() *EventRepo {
	return &EventRepo{
		events: make(map[int64]*model.Event),
		lastId: -1,
	}
}

func (er *EventRepo) Create(event *model.Event) (int64, error) {
	if er.messageAlreadyUsed(event) {
		return -1, ErrEventRepo{
			EventId:  -1,
			Message:  "event message already used",
			Location: "EventRepo.Create()",
		}
	}

	er.lastId++
	event.ID = er.lastId
	er.events[er.lastId] = event
	return er.lastId, nil
}

func (er *EventRepo) Get(id int64) (*model.Event, error) {
	event, ok := er.events[id]
	if !ok {
		return nil, ErrEventRepo{
			EventId:  id,
			Message:  "event not found",
			Location: "EventRepo.Get()",
		}
	}
	return event, nil
}

func (er *EventRepo) GetAll() []*model.Event {
	bonuses := make([]*model.Event, 0, len(er.events))
	for _, bonus := range er.events {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (er *EventRepo) Update(id int64, event *model.Event) error {
	_, ok := er.events[id]
	if !ok {
		return ErrEventRepo{
			EventId:  id,
			Message:  "event not found",
			Location: "EventRepo.Update()",
		}
	}
	if id != event.ID {
		return ErrEventRepo{
			EventId:  id,
			Message:  "events id mismatch",
			Location: "EventRepo.Update()",
		}
	}

	if er.messageAlreadyUsed(event) {
		return ErrEventRepo{
			EventId:  id,
			Message:  "event message already used",
			Location: "EventRepo.Update()",
		}
	}

	er.events[id] = event
	return nil
}

func (er *EventRepo) Delete(id int64) error {
	_, ok := er.events[id]
	if !ok {
		return ErrEventRepo{
			EventId:  id,
			Message:  "event not found",
			Location: "EventRepo.Delete()",
		}
	}
	delete(er.events, id)
	return nil
}

func (er *EventRepo) messageAlreadyUsed(event *model.Event) bool {
	for _, e := range er.events {
		if e.Message == event.Message {
			return true
		}
	}
	return false
}
