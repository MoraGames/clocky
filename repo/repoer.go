package repo

import "github.com/MoraGames/clockyuwu/model"

type UserRepoer interface {
	Create(*model.User) error
	Get(int64) (*model.User, error)
	GetAll() []*model.User
	Update(int64, *model.User) error
	Delete(int64) error
}
