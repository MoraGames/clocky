package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.SetRepoer = new(SetRepo)

// mock.UserRepo
type SetRepo struct {
	sets   map[int64]*model.Set
	lastID int64
}

// Return a new UserRepo
func NewSetRepo() *SetRepo {
	return &SetRepo{
		sets:   make(map[int64]*model.Set),
		lastID: -1,
	}
}

func (sr *SetRepo) Create(set *model.Set) (int64, error) {
	sr.lastID++
	set.ID = sr.lastID
	sr.sets[sr.lastID] = set
	return sr.lastID, nil
}

func (sr *SetRepo) Get(id int64) (*model.Set, error) {
	bonus, ok := sr.sets[id]
	if !ok {
		return nil, errorType.ErrSetNotFound{
			PartecipationID: id,
			Message:         "cannot get set not found",
			Location:        "SetRepo.Get()",
		}
	}
	return bonus, nil
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
		return errorType.ErrSetNotFound{
			SetID:    id,
			Message:  "cannot get set not found",
			Location: "SetRepo.Update()",
		}
	}
	if id != set.ID {
		return errorType.ErrSetNotValid{
			SetID:    id,
			Message:  "cannot update set when id mismatch",
			Location: "SetRepo.Update()",
		}
	}

	sr.sets[id] = set
	return nil
}

func (sr *SetRepo) Delete(id int64) error {
	_, ok := sr.sets[id]
	if !ok {
		return errorType.ErrSetNotFound{
			BonusID:  id,
			Message:  "cannot delete set not found",
			Location: "SetRepo.Delete()",
		}
	}
	delete(sr.sets, id)
	return nil
}
