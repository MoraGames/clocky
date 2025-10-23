package mock

import (
	"fmt"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/repo"
)

// ChatRepo Error
type ErrChatRepo struct {
	ChatId   int64
	Message  string
	Location string
}

func (err ErrChatRepo) Error() string {
	return fmt.Sprintf("%v: %v {id=%v}", err.Location, err.Message, err.ChatId)
}

// Check if the repo implements the interface
var _ repo.ChatRepoer = new(ChatRepo)

// ChatRepo is a mock implementation
type ChatRepo struct {
	chats  map[int64]*model.Chat
	lastId int64
}

func NewChatRepo() *ChatRepo {
	return &ChatRepo{
		chats:  make(map[int64]*model.Chat),
		lastId: -1,
	}
}

func (cr *ChatRepo) Create(chat *model.Chat) (int64, error) {
	cr.lastId++
	chat.ID = cr.lastId
	cr.chats[cr.lastId] = chat
	return cr.lastId, nil
}

func (cr *ChatRepo) Get(id int64) (*model.Chat, error) {
	chat, ok := cr.chats[id]
	if !ok {
		return nil, ErrChatRepo{
			ChatId:   id,
			Message:  "chat not found",
			Location: "ChatRepo.Get()",
		}
	}
	return chat, nil
}

func (cr *ChatRepo) GetAll() []*model.Chat {
	chats := make([]*model.Chat, 0, len(cr.chats))
	for _, chat := range cr.chats {
		chats = append(chats, chat)
	}
	return chats
}

func (cr *ChatRepo) Update(id int64, chat *model.Chat) error {
	_, ok := cr.chats[id]
	if !ok {
		return ErrChatRepo{
			ChatId:   id,
			Message:  "chat not found",
			Location: "ChatRepo.Update()",
		}
	}
	cr.chats[id] = chat
	return nil
}

func (cr *ChatRepo) Delete(id int64) error {
	_, ok := cr.chats[id]
	if !ok {
		return ErrChatRepo{
			ChatId:   id,
			Message:  "chat not found",
			Location: "ChatRepo.Delete()",
		}
	}
	delete(cr.chats, id)
	return nil
}
