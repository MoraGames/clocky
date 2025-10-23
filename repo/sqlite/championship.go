package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// ChampionshipRepo Error
type ErrChampionshipRepo struct {
	ChampionshipId int64
	Message        string
	Location       string
}

func (err ErrChampionshipRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.ChampionshipId)
}

// Check if the repo implements the interface
var _ repo.ChampionshipRepoer = new(ChampionshipRepo)

// ChampionshipRepo is a mock implementation
type ChampionshipRepo struct {
	championships map[int64]*model.Championship
	lastId        int64
}

func NewChampionshipRepo() *ChampionshipRepo {
	return &ChampionshipRepo{
		championships: make(map[int64]*model.Championship),
		lastId:        -1,
	}
}

func (cr *ChampionshipRepo) Create(championship *model.Championship) (int64, error) {
	cr.lastId++
	championship.ID = cr.lastId
	cr.championships[cr.lastId] = championship
	return cr.lastId, nil
}

func (cr *ChampionshipRepo) Get(id int64) (*model.Championship, error) {
	championship, ok := cr.championships[id]
	if !ok {
		return nil, ErrChampionshipRepo{
			ChampionshipId: id,
			Message:        "championship not found",
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
		return ErrChampionshipRepo{
			ChampionshipId: id,
			Message:        "championship not found",
			Location:       "ChampionshipRepo.Update()",
		}
	}
	if id != championship.ID {
		return ErrChampionshipRepo{
			ChampionshipId: id,
			Message:        "championships id mismatch",
			Location:       "ChampionshipRepo.Update()",
		}
	}

	cr.championships[id] = championship
	return nil
}

func (cr *ChampionshipRepo) Delete(id int64) error {
	_, ok := cr.championships[id]
	if !ok {
		return ErrChampionshipRepo{
			ChampionshipId: id,
			Message:        "championship not found",
			Location:       "ChampionshipRepo.Delete()",
		}
	}
	delete(cr.championships, id)
	return nil
}
