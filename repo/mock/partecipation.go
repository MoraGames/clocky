package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.PartecipationRepoer = new(PartecipationRepo)

// mock.UserRepo
type PartecipationRepo struct {
	partecipations map[int64]*model.Partecipation
	lastID         int64
}

// Return a new UserRepo
func NewPartecipationRepo() *PartecipationRepo {
	return &PartecipationRepo{
		partecipations: make(map[int64]*model.Partecipation),
		lastID:         -1,
	}
}

func (pr *PartecipationRepo) Create(partecipation *model.Partecipation) (int64, error) {
	pr.lastID++
	partecipation.ID = pr.lastID
	pr.partecipations[pr.lastID] = partecipation
	return pr.lastID, nil
}

func (pr *PartecipationRepo) Get(id int64) (*model.Partecipation, error) {
	bonus, ok := pr.partecipations[id]
	if !ok {
		return nil, errorType.ErrPartecipationNotFound{
			PartecipationID: id,
			Message:         "cannot get partecipation not found",
			Location:        "PartecipationRepo.Get()",
		}
	}
	return bonus, nil
}

func (pr *PartecipationRepo) GetAll() []*model.Partecipation {
	partecipations := make([]*model.Partecipation, 0, len(pr.partecipations))
	for _, partecipation := range pr.partecipations {
		partecipations = append(partecipations, partecipation)
	}
	return partecipations
}

func (pr *PartecipationRepo) Update(id int64, partecipation *model.Partecipation) error {
	_, ok := pr.partecipations[id]
	if !ok {
		return errorType.ErrPartecipationNotFound{
			PartecipationID: id,
			Message:         "cannot get partecipation not found",
			Location:        "PartecipationRepo.Update()",
		}
	}
	if id != partecipation.ID {
		return errorType.ErrPartecipationNotValid{
			PartecipationID: id,
			Message:         "cannot update partecipation when id mismatch",
			Location:        "PartecipationRepo.Update()",
		}
	}

	pr.partecipations[id] = partecipation
	return nil
}

func (pr *PartecipationRepo) Delete(id int64) error {
	_, ok := pr.partecipations[id]
	if !ok {
		return errorType.ErrPartecipationNotFound{
			PartecipationID: id,
			Message:         "cannot delete partecipation not found",
			Location:        "PartecipationRepo.Delete()",
		}
	}
	delete(pr.partecipations, id)
	return nil
}
