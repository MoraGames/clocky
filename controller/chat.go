package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	"github.com/MoraGames/clockyuwu/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type ChatController struct {
	repo repo.ChatRepoer
	log  *logrus.Logger
}

func NewChatController(repoer repo.ChatRepoer, logger *logrus.Logger) *ChatController {
	return &ChatController{
		repo: repoer,
		log:  logger,
	}
}

func (cc *ChatController) CreateChat(chatID int64, telegramChat *tgbotapi.Chat) error {
	//Check if the user already exists
	if _, err := cc.repo.Get(chatID); err == nil {
		return errorType.ErrChatAlreadyExist{
			ChatID:   chatID,
			Message:  "cannot create chat that already exists",
			Location: "ChatController.CreateChat()",
		}
	} else if err.Error() != "cannot get chat not found" {
		return err
	}

	//Create the user
	chat := &model.Chat{
		TelegramChat:  telegramChat,
		Type:          telegramChat.Type,
		Title:         telegramChat.Title,
		Championships: make([]*model.Championship, 0),
	}

	return cc.repo.Create(chat)
}

func (cc *ChatController) GetChat(chatID int64) (*model.Chat, error) {
	return cc.repo.Get(chatID)
}

func (cc *ChatController) GetAllChats() []*model.Chat {
	return cc.repo.GetAll()
}

func (cc *ChatController) GetChatChampionships(chatID int64) ([]*model.Championship, error) {
	//Check if the chat already exists
	chat, err := cc.repo.Get(chatID)
	if err != nil {
		return nil, err
	}

	return chat.Championships, nil
}

func (cc *ChatController) DeleteChat(chatID int64) error {
	//Check if the chat already exists
	_, err := cc.repo.Get(chatID)
	if err != nil {
		return err
	}

	//Delete the chat
	return cc.repo.Delete(chatID)
}
