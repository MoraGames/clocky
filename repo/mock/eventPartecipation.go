package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// EventPartecipationRepo Error
type ErrEventPartecipationRepo struct {
	EventPartecipationId int64
	Message              string
	Location             string
}

func (err ErrEventPartecipationRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.EventPartecipationId)
}

// Check if the repo implements the interface
var _ repo.EventPartecipationRepoer = new(EventPartecipationRepo)

// EventPartecipationRepo is a mock implementation
type EventPartecipationRepo struct {
	eventPartecipations map[int64]*model.EventPartecipation
	lastId              int64
}

func NewEventPartecipationRepo() *EventPartecipationRepo {
	return &EventPartecipationRepo{
		eventPartecipations: make(map[int64]*model.EventPartecipation),
		lastId:              -1,
	}
}

func (epr *EventPartecipationRepo) Create(eventPartecipation *model.EventPartecipation) (int64, error) {
	epr.lastId++
	eventPartecipation.ID = epr.lastId
	epr.eventPartecipations[epr.lastId] = eventPartecipation
	return epr.lastId, nil
}

func (epr *EventPartecipationRepo) Get(id int64) (*model.EventPartecipation, error) {
	eventPartecipation, ok := epr.eventPartecipations[id]
	if !ok {
		return nil, ErrEventPartecipationRepo{
			EventPartecipationId: id,
			Message:              "eventPartecipation not found",
			Location:             "EventPartecipationRepo.Get()",
		}
	}
	return eventPartecipation, nil
}

func (epr *EventPartecipationRepo) GetAll() []*model.EventPartecipation {
	bonuses := make([]*model.EventPartecipation, 0, len(epr.eventPartecipations))
	for _, bonus := range epr.eventPartecipations {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (epr *EventPartecipationRepo) Update(id int64, eventPartecipation *model.EventPartecipation) error {
	_, ok := epr.eventPartecipations[id]
	if !ok {
		return ErrEventPartecipationRepo{
			EventPartecipationId: id,
			Message:              "eventPartecipation not found",
			Location:             "EventPartecipationRepo.Update()",
		}
	}
	if id != eventPartecipation.ID {
		return ErrEventPartecipationRepo{
			EventPartecipationId: id,
			Message:              "eventPartecipations id mismatch",
			Location:             "EventPartecipationRepo.Update()",
		}
	}

	epr.eventPartecipations[id] = eventPartecipation
	return nil
}

func (epr *EventPartecipationRepo) Delete(id int64) error {
	_, ok := epr.eventPartecipations[id]
	if !ok {
		return ErrEventPartecipationRepo{
			EventPartecipationId: id,
			Message:              "eventPartecipation not found",
			Location:             "EventPartecipationRepo.Delete()",
		}
	}
	delete(epr.eventPartecipations, id)
	return nil
}
