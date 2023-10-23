package repo

import "github.com/MoraGames/clockyuwu/model"

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

type RecordRepoer interface {
	Create(*model.Record) error
	Get(string) (*model.Record, error)
	GetAll() []*model.Record
	Update(string, *model.Record) error
	Delete(string) error
}

type UserRepoer interface {
	Create(*model.User) error
	Get(int64) (*model.User, error)
	GetAll() []*model.User
	Update(int64, *model.User) error
	Delete(int64) error
}
