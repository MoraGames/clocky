package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/MoraGames/clockyuwu/controller"
	"github.com/MoraGames/clockyuwu/pkg/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Command struct {
	Command      string
	Arguments    []string
	Description  string
	InputScopes  []string
	OutputScopes []string
	AdminOnly    bool
}

func manageCommand(appUtils util.AppUtils, controller *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time) {
	switch update.Message.Command() {
	case "credits":
		cmd := Command{
			Command:      "/credits",
			Arguments:    strings.Split(update.Message.CommandArguments(), " "),
			Description:  "Get more informations abount the project.",
			InputScopes:  []string{"private", "group", "supergroup"},
			OutputScopes: []string{"private"},
			AdminOnly:    false,
		}
		cmdCredits(appUtils, controller, bot, update, updateCurTime, cmd)
	case "help":
		cmd := Command{
			Command:      "/help",
			Arguments:    strings.Split(update.Message.CommandArguments(), " "),
			Description:  "Get a complete list of all available commands.",
			InputScopes:  []string{"private", "group", "supergroup"},
			OutputScopes: []string{"private"},
			AdminOnly:    false,
		}
		cmdHelp(appUtils, controller, bot, update, updateCurTime, cmd)
	case "ping":
		cmd := Command{
			Command:      "/ping",
			Arguments:    strings.Split(update.Message.CommandArguments(), " "),
			Description:  "Verify if the bot is running.",
			InputScopes:  []string{"private", "group", "supergroup"},
			OutputScopes: []string{"private", "group", "supergroup"},
			AdminOnly:    false,
		}
		cmdPing(appUtils, controller, bot, update, updateCurTime, cmd)
	case "start":
		cmd := Command{
			Command:      "/start",
			Arguments:    strings.Split(update.Message.CommandArguments(), " "),
			Description:  "Get an introductory message about the bot's features.",
			InputScopes:  []string{"private", "group", "supergroup"},
			OutputScopes: []string{"private"},
			AdminOnly:    false,
		}
		cmdStart(appUtils, controller, bot, update, updateCurTime, cmd)
	default:
		// Respond with an error message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.\nUse /help to get a list of all commands.")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		appUtils.Logger.WithFields(logrus.Fields{
			"message": update.Message.Text,
			"sender":  update.Message.From.UserName,
			"chat":    update.Message.Chat.Title,
		}).Debug("Response to unknown command sent successfully")
	}
}

func cmdCredits(appUtils util.AppUtils, controller *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, cmd Command) {
	// Respond with useful information about the project
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The source code, available on GitHub at MoraGames/clockyuwu, is written entirely in GoLang and makes use of the \"telegram-bot-api\" library.\nFor any bug reports or feature proposals, please refer to the GitHub project.\n\nDeveloper:\n- Telegram: @MoraGames\n- Discord: @moragames\n- Instagram: @moragames.dev\n- GitHub: MoraGames\n\nProject:\n- Telegram: @clockyuwu_bot\n- GitHub: MoraGames/clockyuwu\n\nSpecial thanks go to the first testers (as well as players) of the minigame managed by the bot, \"Vano\", \"Ale\" and \"Alex\".")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	appUtils.Logger.WithFields(logrus.Fields{
		"message": update.Message.Text,
		"sender":  update.Message.From.UserName,
		"chat":    update.Message.Chat.Title,
	}).Debug("Response to \"/credits\" command sent successfully")
}

func cmdHelp(appUtils util.AppUtils, controller *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, cmd Command) {
	// Respond with useful information about the working and commands of the bot
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Name: %v\nVersion: %v\n\nThis is a list of all possible commands within the bot:\n\n- /start : Get an introductory message about the bot's features.\n - /help : Get a complete list of all available commands.\n - /ping : Verify if the bot is running.\n - /credits : Get more informations abount the project.", appUtils.ConfigApp.Name, appUtils.ConfigApp.Version))
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	appUtils.Logger.WithFields(logrus.Fields{
		"message": update.Message.Text,
		"sender":  update.Message.From.UserName,
		"chat":    update.Message.Chat.Title,
	}).Debug("Response to \"/help\" command sent successfully")
}

func cmdPing(appUtils util.AppUtils, controller *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, cmd Command) {
	// Respond with a "pong" message. Useful for checking if the bot is online
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	appUtils.Logger.WithFields(logrus.Fields{
		"message": update.Message.Text,
		"sender":  update.Message.From.UserName,
		"chat":    update.Message.Chat.Title,
	}).Debug("Response to \"/ping\" command sent successfully")
}

func cmdStart(appUtils util.AppUtils, controller *controller.Controller, bot *tgbotapi.BotAPI, update tgbotapi.Update, updateCurTime time.Time, cmd Command) {
	// Respond with an introduction message for the users of the bot
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%v is a bot that allows you to play a time-wasting game with one or more groups of friends within Telegram groups. Once the bot is added, the game mainly (but not exclusively) involves sending messages in the \"hh:mm\" format at certain times of the day, in exchange for valuable points. The person who has earned the most points at the end of the championship will be the new Clocky Champion!\nUse /help to get a list of all commands or /credits for more information about the project.\n\n- %v, a bot from @MoraGames.", appUtils.ConfigApp.Name, appUtils.ConfigApp.Name))
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	appUtils.Logger.WithFields(logrus.Fields{
		"message": update.Message.Text,
		"sender":  update.Message.From.UserName,
		"chat":    update.Message.Chat.Title,
	}).Debug("Response to \"/start\" command sent successfully")
}
