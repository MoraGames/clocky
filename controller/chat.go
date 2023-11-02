package controller

import (
	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/errorType"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Controller) CreateChat(chatID int64, telegramChat *tgbotapi.Chat) error {
	//Check if the user already exists
	if _, err := c.chat.Get(chatID); err == nil {
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

	return c.chat.Create(chat)
}

func (c *Controller) GetChat(chatID int64) (*model.Chat, error) {
	return c.chat.Get(chatID)
}

func (c *Controller) GetAllChats() []*model.Chat {
	return c.chat.GetAll()
}

func (c *Controller) GetChatChampionships(chatID int64) ([]*model.Championship, error) {
	//Check if the chat already exists
	chat, err := c.chat.Get(chatID)
	if err != nil {
		return nil, err
	}

	return chat.Championships, nil
}

func (c *Controller) DeleteChat(chatID int64) error {
	//Check if the chat already exists
	_, err := c.chat.Get(chatID)
	if err != nil {
		return err
	}

	//Delete the chat
	return c.chat.Delete(chatID)
}
