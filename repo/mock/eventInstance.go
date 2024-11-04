package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// EventInstanceRepo Error
type ErrEventInstanceRepo struct {
	EventInstanceId int64
	Message         string
	Location        string
}

func (err ErrEventInstanceRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.EventInstanceId)
}

// Check if the repo implements the interface
var _ repo.EventInstanceRepoer = new(EventInstanceRepo)

// EventInstanceRepo is a mock implementation
type EventInstanceRepo struct {
	eventInstances map[int64]*model.EventInstance
	lastId         int64
}

func NewEventInstanceRepo() *EventInstanceRepo {
	return &EventInstanceRepo{
		eventInstances: make(map[int64]*model.EventInstance),
		lastId:         -1,
	}
}

func (eir *EventInstanceRepo) Create(eventInstance *model.EventInstance) (int64, error) {
	eir.lastId++
	eventInstance.ID = eir.lastId
	eir.eventInstances[eir.lastId] = eventInstance
	return eir.lastId, nil
}

func (eir *EventInstanceRepo) Get(id int64) (*model.EventInstance, error) {
	eventInstance, ok := eir.eventInstances[id]
	if !ok {
		return nil, ErrEventInstanceRepo{
			EventInstanceId: id,
			Message:         "eventInstance not found",
			Location:        "EventInstanceRepo.Get()",
		}
	}
	return eventInstance, nil
}

func (eir *EventInstanceRepo) GetAll() []*model.EventInstance {
	bonuses := make([]*model.EventInstance, 0, len(eir.eventInstances))
	for _, bonus := range eir.eventInstances {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (eir *EventInstanceRepo) Update(id int64, eventInstance *model.EventInstance) error {
	_, ok := eir.eventInstances[id]
	if !ok {
		return ErrEventInstanceRepo{
			EventInstanceId: id,
			Message:         "eventInstance not found",
			Location:        "EventInstanceRepo.Update()",
		}
	}
	if id != eventInstance.ID {
		return ErrEventInstanceRepo{
			EventInstanceId: id,
			Message:         "eventInstances id mismatch",
			Location:        "EventInstanceRepo.Update()",
		}
	}

	eir.eventInstances[id] = eventInstance
	return nil
}

func (eir *EventInstanceRepo) Delete(id int64) error {
	_, ok := eir.eventInstances[id]
	if !ok {
		return ErrEventInstanceRepo{
			EventInstanceId: id,
			Message:         "eventInstance not found",
			Location:        "EventInstanceRepo.Delete()",
		}
	}
	delete(eir.eventInstances, id)
	return nil
}

func (eir *EventInstanceRepo) GetLastEventInstanceByEvent(event *model.Event) (*model.EventInstance, error) {
	for _, eventInstance := range eir.eventInstances {
		if eventInstance.Event.ID == event.ID {
			return eventInstance, nil
		}
	}
	return nil, ErrEventInstanceRepo{
		EventInstanceId: event.ID,
		Message:         "eventInstance not found",
		Location:        "EventInstanceRepo.GetLastEventInstanceByEvent()",
	}
}
