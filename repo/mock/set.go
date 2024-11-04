package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// SetRepo Error
type ErrSetRepo struct {
	SetId    int64
	Message  string
	Location string
}

func (err ErrSetRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.SetId)
}

// Check if the repo implements the interface
var _ repo.SetRepoer = new(SetRepo)

// SetRepo is a mock implementation
type SetRepo struct {
	sets   map[int64]*model.Set
	lastId int64
}

func NewSetRepo() *SetRepo {
	return &SetRepo{
		sets:   make(map[int64]*model.Set),
		lastId: -1,
	}
}

func (sr *SetRepo) Create(set *model.Set) (int64, error) {
	sr.lastId++
	set.ID = sr.lastId
	sr.sets[sr.lastId] = set
	return sr.lastId, nil
}

func (sr *SetRepo) Get(id int64) (*model.Set, error) {
	set, ok := sr.sets[id]
	if !ok {
		return nil, ErrSetRepo{
			SetId:    id,
			Message:  "set not found",
			Location: "SetRepo.Get()",
		}
	}
	return set, nil
}

func (sr *SetRepo) GetAll() []*model.Set {
	sets := make([]*model.Set, 0, len(sr.sets))
	for _, set := range sr.sets {
		sets = append(sets, set)
	}
	return sets
}

func (sr *SetRepo) Update(id int64, set *model.Set) error {
	_, ok := sr.sets[id]
	if !ok {
		return ErrSetRepo{
			SetId:    id,
			Message:  "set not found",
			Location: "SetRepo.Update()",
		}
	}
	if id != set.ID {
		return ErrSetRepo{
			SetId:    id,
			Message:  "sets id mismatch",
			Location: "SetRepo.Update()",
		}
	}

	sr.sets[id] = set
	return nil
}

func (sr *SetRepo) Delete(id int64) error {
	_, ok := sr.sets[id]
	if !ok {
		return ErrSetRepo{
			SetId:    id,
			Message:  "set not found",
			Location: "SetRepo.Delete()",
		}
	}
	delete(sr.sets, id)
	return nil
}
