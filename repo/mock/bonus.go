package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.BonusRepoer = new(BonusRepo)

// mock.UserRepo
type BonusRepo struct {
	bonuses map[int64]*model.Bonus
	lastID  int64
}

// Return a new UserRepo
func NewBonusRepo() *BonusRepo {
	return &BonusRepo{
		bonuses: make(map[int64]*model.Bonus),
		lastID:  -1,
	}
}

func (br *BonusRepo) Create(bonus *model.Bonus) (int64, error) {
	br.lastID++
	bonus.ID = br.lastID
	br.bonuses[br.lastID] = bonus
	return br.lastID, nil
}

func (br *BonusRepo) Get(id int64) (*model.Bonus, error) {
	bonus, ok := br.bonuses[id]
	if !ok {
		return nil, errorType.ErrBonusNotFound{
			BonusID:  id,
			Message:  "cannot get bonus not found",
			Location: "BonusRepo.Get()",
		}
	}
	return bonus, nil
}

func (br *BonusRepo) GetAll() []*model.Bonus {
	bonuses := make([]*model.Bonus, 0, len(br.bonuses))
	for _, bonus := range br.bonuses {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (br *BonusRepo) Update(id int64, bonus *model.Bonus) error {
	_, ok := br.bonuses[id]
	if !ok {
		return errorType.ErrBonusNotFound{
			BonusID:  id,
			Message:  "cannot get bonus not found",
			Location: "BonusRepo.Update()",
		}
	}
	if id != bonus.ID {
		return errorType.ErrBonusNotValid{
			BonusID:  id,
			Message:  "cannot update bonus when id mismatch",
			Location: "BonusRepo.Update()",
		}
	}

	br.bonuses[id] = bonus
	return nil
}

func (br *BonusRepo) Delete(id int64) error {
	_, ok := br.bonuses[id]
	if !ok {
		return errorType.ErrBonusNotFound{
			BonusID:  id,
			Message:  "cannot delete bonus not found",
			Location: "BonusRepo.Delete()",
		}
	}
	delete(br.bonuses, id)
	return nil
}
