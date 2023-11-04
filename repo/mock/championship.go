package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.ChampionshipRepoer = new(ChampionshipRepo)

// mock.UserRepo
type ChampionshipRepo struct {
	championships map[int64]*model.Championship
	lastID        int64
}

// Return a new UserRepo
func NewChampionshipRepo() *ChampionshipRepo {
	return &ChampionshipRepo{
		championships: make(map[int64]*model.Championship),
		lastID:        -1,
	}
}

func (cr *ChampionshipRepo) Create(championship *model.Championship) (int64, error) {
	cr.lastID++
	championship.ID = cr.lastID
	cr.championships[cr.lastID] = championship
	return cr.lastID, nil
}

func (cr *ChampionshipRepo) Get(id int64) (*model.Championship, error) {
	championship, ok := cr.championships[id]
	if !ok {
		return nil, errorType.ErrChampionshipNotFound{
			ChampionshipID: id,
			Message:        "cannot get championship not found",
			Location:       "ChampionshipRepo.Get()",
		}
	}
	return championship, nil
}

func (cr *ChampionshipRepo) GetAll() []*model.Championship {
	championships := make([]*model.Championship, 0, len(cr.championships))
	for _, championship := range cr.championships {
		championships = append(championships, championship)
	}
	return championships
}

func (cr *ChampionshipRepo) Update(id int64, championship *model.Championship) error {
	_, ok := cr.championships[id]
	if !ok {
		return errorType.ErrChampionshipNotFound{
			ChampionshipID: id,
			Message:        "cannot update championship not found",
			Location:       "ChampionshipRepo.Update()",
		}
	}
	if id != championship.ID {
		return errorType.ErrChampionshipNotValid{
			ChampionshipID: id,
			Message:        "cannot update championship when id mismatch",
			Location:       "ChampionshipRepo.Update()",
		}
	}

	cr.championships[id] = championship
	return nil
}

func (cr *ChampionshipRepo) Delete(id int64) error {
	_, ok := cr.championships[id]
	if !ok {
		return errorType.ErrChampionshipNotFound{
			ChampionshipID: id,
			Message:        "cannot delete championship not found",
			Location:       "ChampionshipRepo.Delete()",
		}
	}
	delete(cr.championships, id)
	return nil
}
