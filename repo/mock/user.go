package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.UserRepoer = new(UserRepo)

// mock.UserRepo
type UserRepo struct {
	users map[int64]*model.User
}

// Return a new UserRepo
func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[int64]*model.User),
	}
}

func (ur *UserRepo) Create(user *model.User) error {
	ur.users[user.TelegramUser.ID] = user
	return nil
}

func (ur *UserRepo) Get(id int64) (*model.User, error) {
	user, ok := ur.users[id]
	if !ok {
		return nil, errorType.ErrUserNotFound{
			UserID:   id,
			Message:  "cannot get user not found",
			Location: "UserRepo.Get()",
		}
	}
	return user, nil
}

func (ur *UserRepo) GetAll() ([]*model.User, error) {
	users := make([]*model.User, 0, len(ur.users))
	for _, user := range ur.users {
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRepo) Update(id int64, user *model.User) error {
	_, ok := ur.users[id]
	if !ok {
		return errorType.ErrUserNotFound{
			UserID:   id,
			Message:  "cannot update user not found",
			Location: "UserRepo.Update()",
		}
	}
	ur.users[id] = user
	return nil
}

func (ur *UserRepo) Delete(id int64) error {
	_, ok := ur.users[id]
	if !ok {
		return errorType.ErrUserNotFound{
			UserID:   id,
			Message:  "cannot delete user not found",
			Location: "UserRepo.Delete()",
		}
	}
	delete(ur.users, id)
	return nil
}
