package repo

import "github.com/MoraGames/clockyuwu/model"

type BonusRepoer interface {
	Create(*model.Bonus) (int64, error)
	Get(int64) (*model.Bonus, error)
	GetAll() []*model.Bonus
	Update(int64, *model.Bonus) error
	Delete(int64) error
}
type ChampionshipRepoer interface {
	Create(*model.Championship) (int64, error)
	Get(int64) (*model.Championship, error)
	GetAll() []*model.Championship
	GetLast() (*model.Championship, error)
	Update(int64, *model.Championship) error
	Delete(int64) error
}
type ChatRepoer interface {
	Create(*model.Chat) error
	Get(int64) (*model.Chat, error)
	GetAll() []*model.Chat
	Update(int64, *model.Chat) error
	Delete(int64) error
}
type EventRepoer interface {
	Create(*model.Event) error
	Get(string) (*model.Event, error)
	GetAll() []*model.Event
	Update(string, *model.Event) error
	Delete(string) error
}
type PartecipationRepoer interface {
	Create(*model.Partecipation) (int64, error)
	Get(int64) (*model.Partecipation, error)
	GetAll() []*model.Partecipation
	Update(int64, *model.Partecipation) error
	Delete(int64) error
}
type RecordRepoer interface {
	Create(*model.Record) error
	Get(string) (*model.Record, error)
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
	Create(*model.User) error
	Get(int64) (*model.User, error)
	GetAll() []*model.User
	Update(int64, *model.User) error
	Delete(int64) error
}
