package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.EffectRepoer = new(EffectRepo)

// mock.UserRepo
type EffectRepo struct {
	effects map[string]*model.Effect
}

// Return a new UserRepo
func NewBonusRepo() *EffectRepo {
	return &EffectRepo{
		effects: make(map[string]*model.Effect),
	}
}

func (er *EffectRepo) Create(effect *model.Effect) error {
	if _, ok := er.effects[effect.Name]; ok {
		return errorType.ErrEffectAlreadyExist{
			EffectName: effect.Name,
			Message:    "cannot create effect that already exists",
			Location:   "EffectRepo.Create()",
		}
	}

	er.effects[effect.Name] = effect
	return nil
}

func (er *EffectRepo) Get(name string) (*model.Effect, error) {
	bonus, ok := er.effects[name]
	if !ok {
		return nil, errorType.ErrEffectNotFound{
			EffectName: name,
			Message:    "cannot get effect not found",
			Location:   "EffectRepo.Get()",
		}
	}
	return bonus, nil
}

func (er *EffectRepo) GetAll() []*model.Effect {
	bonuses := make([]*model.Effect, 0, len(er.effects))
	for _, bonus := range er.effects {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (er *EffectRepo) Update(name string, effect *model.Effect) error {
	_, ok := er.effects[name]
	if !ok {
		return errorType.ErrEffectNotFound{
			EffectName: name,
			Message:    "cannot get effect not found",
			Location:   "EffectRepo.Update()",
		}
	}
	if name != effect.Name {
		return errorType.ErrEffectNotValid{
			EffectName: name,
			Message:    "cannot update effect when id mismatch",
			Location:   "EffectRepo.Update()",
		}
	}

	er.effects[name] = effect
	return nil
}

func (er *EffectRepo) Delete(name string) error {
	_, ok := er.effects[name]
	if !ok {
		return errorType.ErrEffectNotFound{
			EffectName: name,
			Message:    "cannot delete effect not found",
			Location:   "EffectRepo.Delete()",
		}
	}
	delete(er.effects, name)
	return nil
}
