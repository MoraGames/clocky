package mock

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
)

// Check if the repo implements the interface
var _ repo.ChatRepoer = new(ChatRepo)

// mock.UserRepo
type ChatRepo struct {
	chats map[int64]*model.Chat
}

// Return a new UserRepo
func NewChatRepo() *ChatRepo {
	return &ChatRepo{
		chats: make(map[int64]*model.Chat),
	}
}

func (cr *ChatRepo) Create(chat *model.Chat) error {
	cr.chats[chat.TelegramChat.ID] = chat
	return nil
}

func (cr *ChatRepo) Get(id int64) (*model.Chat, error) {
	chat, ok := cr.chats[id]
	if !ok {
		return nil, errorType.ErrChatNotFound{
			ChatID:   id,
			Message:  "cannot get chat not found",
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
		return errorType.ErrChatNotFound{
			ChatID:   id,
			Message:  "cannot update chat not found",
			Location: "ChatRepo.Update()",
		}
	}
	cr.chats[id] = chat
	return nil
}

func (cr *ChatRepo) Delete(id int64) error {
	_, ok := cr.chats[id]
	if !ok {
		return errorType.ErrChatNotFound{
			ChatID:   id,
			Message:  "cannot delete chat not found",
			Location: "ChatRepo.Delete()",
		}
	}
	delete(cr.chats, id)
	return nil
}
