package repo

import "github.com/MoraGames/clockyuwu/model"

type ChampionshipRepoer interface {
	Create(*model.Championship) (int64, error)
	Get(int64) (*model.Championship, error)
	GetAll() []*model.Championship
	Update(int64, *model.Championship) error
	Delete(int64) error
}
type ChatRepoer interface {
	Create(*model.Chat) (int64, error)
	Get(int64) (*model.Chat, error)
	GetAll() []*model.Chat
	Update(int64, *model.Chat) error
	Delete(int64) error
}
type EffectRepoer interface {
	Create(*model.Effect) (int64, error)
	Get(int64) (*model.Effect, error)
	GetAll() []*model.Effect
	Update(string, *model.Effect) error
	Delete(string) error
}
type EventRepoer interface {
	Create(*model.Event) (int64, error)
	Get(int64) (*model.Event, error)
	GetAll() []*model.Event
	Update(string, *model.Event) error
	Delete(string) error
}
type EventInstanceReport interface {
	Create(*model.EventInstance) (int64, error)
	Get(int64) (*model.EventInstance, error)
	GetAll() []*model.EventInstance
	Update(int64, *model.EventInstance) error
	Delete(int64) error
}
type EventPartecipationRepoer interface {
	Create(*model.EventPartecipation) (int64, error)
	Get(int64) (*model.EventPartecipation, error)
	GetAll() []*model.EventPartecipation
	Update(int64, *model.EventPartecipation) error
	Delete(int64) error
}
type RecordRepoer interface {
	Create(*model.Record) (int64, error)
	Get(int64) (*model.Record, error)
	GetAll() []*model.Record
	Update(string, *model.Record) error
	Delete(string) error
}
type SetRepoer interface {
	Create(*model.Set) (int64, error)
	Get(int64) (*model.Set, error)
	GetAll() []*model.Set
	Update(int64, *model.Set) error
	Delete(int64) error
}
type UserRepoer interface {
	Create(*model.User) (int64, error)
	Get(int64) (*model.User, error)
	GetAll() []*model.User
	Update(int64, *model.User) error
	Delete(int64) error
}
