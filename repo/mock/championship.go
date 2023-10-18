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
	lastEdition   int64
}

// Return a new UserRepo
func NewChampionshipRepo() *ChampionshipRepo {
	return &ChampionshipRepo{
		championships: make(map[int64]*model.Championship),
		lastEdition:   0,
	}
}

func (cr *ChampionshipRepo) Create(championship *model.Championship) error {
	cr.championships[cr.lastEdition+1] = championship
	cr.lastEdition++
	return nil
}

func (cr *ChampionshipRepo) Get(edition int64) (*model.Championship, error) {
	user, ok := cr.championships[edition]
	if !ok {
		return nil, errorType.ErrChampionshipNotFound{
			ChampionshipEdition: edition,
			Message:             "cannot get championship not found",
			Location:            "ChampionshipRepo.Get()",
		}
	}
	return user, nil
}

func (cr *ChampionshipRepo) GetAll() []*model.Championship {
	championships := make([]*model.Championship, 0, len(cr.championships))
	for _, championship := range cr.championships {
		championships = append(championships, championship)
	}
	return championships
}

func (cr *ChampionshipRepo) GetLast() (*model.Championship, error) {
	user, ok := cr.championships[cr.lastEdition]
	if !ok {
		return nil, errorType.ErrChampionshipNotFound{
			ChampionshipEdition: cr.lastEdition,
			Message:             "cannot get championship not found",
			Location:            "ChampionshipRepo.GetLast()",
		}
	}
	return user, nil
}

func (cr *ChampionshipRepo) Update(edition int64, championship *model.Championship) error {
	_, ok := cr.championships[edition]
	if !ok {
		return errorType.ErrChampionshipNotFound{
			ChampionshipEdition: edition,
			Message:             "cannot update championship not found",
			Location:            "ChampionshipRepo.Update()",
		}
	}
	cr.championships[edition] = championship
	return nil
}

func (cr *ChampionshipRepo) Delete(edition int64) error {
	_, ok := cr.championships[edition]
	if !ok {
		return errorType.ErrChampionshipNotFound{
			ChampionshipEdition: edition,
			Message:             "cannot delete championship not found",
			Location:            "ChampionshipRepo.Delete()",
		}
	}
	delete(cr.championships, edition)
	return nil
}
