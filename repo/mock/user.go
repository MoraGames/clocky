package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// UserRepo Error
type ErrUserRepo struct {
	UserId   int64
	Message  string
	Location string
}

func (err ErrUserRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.UserId)
}

// Check if the repo implements the interface
var _ repo.UserRepoer = new(UserRepo)

// UserRepo is a mock implementation
type UserRepo struct {
	users  map[int64]*model.User
	lastId int64
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:  make(map[int64]*model.User),
		lastId: -1,
	}
}

func (ur *UserRepo) Create(user *model.User) (int64, error) {
	if ur.nicknameAlreadyUsed(user) {
		return -1, ErrUserRepo{
			UserId:   -1,
			Message:  "user nickname already used",
			Location: "UserRepo.Create()",
		}
	}

	ur.lastId++
	user.ID = ur.lastId
	ur.users[ur.lastId] = user
	return ur.lastId, nil
}

func (ur *UserRepo) Get(id int64) (*model.User, error) {
	user, ok := ur.users[id]
	if !ok {
		return nil, ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Get()",
		}
	}
	return user, nil
}

func (ur *UserRepo) GetAll() []*model.User {
	users := make([]*model.User, 0, len(ur.users))
	for _, user := range ur.users {
		users = append(users, user)
	}
	return users
}

func (ur *UserRepo) Update(id int64, user *model.User) error {
	_, ok := ur.users[id]
	if !ok {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Update()",
		}
	}
	if id != user.ID {
		return ErrUserRepo{
			UserId:   id,
			Message:  "users id mismatch",
			Location: "UserRepo.Update()",
		}
	}

	if ur.nicknameAlreadyUsed(user) {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user nickname already used",
			Location: "UserRepo.Update()",
		}
	}

	ur.users[id] = user
	return nil
}

func (ur *UserRepo) Delete(id int64) error {
	_, ok := ur.users[id]
	if !ok {
		return ErrUserRepo{
			UserId:   id,
			Message:  "user not found",
			Location: "UserRepo.Delete()",
		}
	}
	delete(ur.users, id)
	return nil
}

func (ur *UserRepo) nicknameAlreadyUsed(user *model.User) bool {
	for _, u := range ur.users {
		if u.Nickname == user.Nickname {
			return true
		}
	}
	return false
}
