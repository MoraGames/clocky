package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// EffectRepo Error
type ErrEffectRepo struct {
	EffectId int64
	Message  string
	Location string
}

func (err ErrEffectRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.EffectId)
}

// Check if the repo implements the interface
var _ repo.EffectRepoer = new(EffectRepo)

// EffectRepo is a mock implementation
type EffectRepo struct {
	effects map[int64]*model.Effect
	lastId  int64
}

func NewEffectRepo() *EffectRepo {
	return &EffectRepo{
		effects: make(map[int64]*model.Effect),
		lastId:  -1,
	}
}

func (er *EffectRepo) Create(effect *model.Effect) (int64, error) {
	if er.nameAlreadyUsed(effect) {
		return -1, ErrEffectRepo{
			EffectId: -1,
			Message:  "effect name already used",
			Location: "EffectRepo.Create()",
		}
	}

	er.lastId++
	effect.ID = er.lastId
	er.effects[er.lastId] = effect
	return er.lastId, nil
}

func (er *EffectRepo) Get(id int64) (*model.Effect, error) {
	effect, ok := er.effects[id]
	if !ok {
		return nil, ErrEffectRepo{
			EffectId: id,
			Message:  "effect not found",
			Location: "EffectRepo.Get()",
		}
	}
	return effect, nil
}

func (er *EffectRepo) GetAll() []*model.Effect {
	bonuses := make([]*model.Effect, 0, len(er.effects))
	for _, bonus := range er.effects {
		bonuses = append(bonuses, bonus)
	}
	return bonuses
}

func (er *EffectRepo) Update(id int64, effect *model.Effect) error {
	_, ok := er.effects[id]
	if !ok {
		return ErrEffectRepo{
			EffectId: id,
			Message:  "effect not found",
			Location: "EffectRepo.Update()",
		}
	}
	if id != effect.ID {
		return ErrEffectRepo{
			EffectId: id,
			Message:  "effects id mismatch",
			Location: "EffectRepo.Update()",
		}
	}

	if er.nameAlreadyUsed(effect) {
		return ErrEffectRepo{
			EffectId: id,
			Message:  "effect name already used",
			Location: "EffectRepo.Update()",
		}
	}

	er.effects[id] = effect
	return nil
}

func (er *EffectRepo) Delete(id int64) error {
	_, ok := er.effects[id]
	if !ok {
		return ErrEffectRepo{
			EffectId: id,
			Message:  "effect not found",
			Location: "EffectRepo.Delete()",
		}
	}
	delete(er.effects, id)
	return nil
}

func (er *EffectRepo) nameAlreadyUsed(effect *model.Effect) bool {
	for _, e := range er.effects {
		if e.Name == effect.Name {
			return true
		}
	}
	return false
}
