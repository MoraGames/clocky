package main

import (
	"fmt"
	"strings"

	"github.com/MoraGames/clockyuwu/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type ShortCommand struct {
	Command      string
	Description  string
	InputScopes  []string
	OutputScopes []string
	AdminOnly    bool
	StaffOnly    bool
}

func (scmd ShortCommand) ShortFormat() string {
	return fmt.Sprintf(" - %v : %v", scmd.Command, scmd.Description)
}

func (scmd ShortCommand) Format() string {
	return fmt.Sprintf(" - %v :\nDescription: %v\nInput: %v\nOutput: %v\nAdmin Only: %v\nStaff Only: %v", scmd.Command, scmd.Description, scmd.InputScopes, scmd.OutputScopes, scmd.AdminOnly, scmd.StaffOnly)
}

func (scmd ShortCommand) ToCommand(commandUpdate utils.CommandUpdate, message string) Command {
	return Command{
		CommandUpdate: commandUpdate,
		Command:       scmd.Command,
		Arguments:     strings.Split(commandUpdate.Update.Message.CommandArguments(), " "),
		Description:   scmd.Description,
		Message:       message,
		InputScopes:   scmd.InputScopes,
		OutputScopes:  scmd.OutputScopes,
		AdminOnly:     scmd.AdminOnly,
		StaffOnly:     scmd.StaffOnly,
	}
}

var ShortCommands = map[string]ShortCommand{
	"credits": {
		Command:      "/credits",
		Description:  "Get more informations abount the project.",
		InputScopes:  []string{"private", "group", "supergroup"},
		OutputScopes: []string{"private"},
		AdminOnly:    false,
		StaffOnly:    false,
	},
	"help": {
		Command:      "/help",
		Description:  "Get a complete list of all available commands.",
		InputScopes:  []string{"private", "group", "supergroup"},
		OutputScopes: []string{"private"},
		AdminOnly:    false,
		StaffOnly:    false,
	},
	"ping": {
		Command:      "/ping",
		Description:  "Verify if the bot is running.",
		InputScopes:  []string{"private", "group", "supergroup"},
		OutputScopes: []string{"private", "group", "supergroup"},
		AdminOnly:    false,
		StaffOnly:    false,
	},
	"start": {
		Command:      "/start",
		Description:  "Get an introductory message about the bot's features.",
		InputScopes:  []string{"private", "group", "supergroup"},
		OutputScopes: []string{"private"},
		AdminOnly:    false,
		StaffOnly:    false,
	},
}

type Command struct {
	CommandUpdate utils.CommandUpdate
	Command       string
	Arguments     []string
	Description   string
	Message       string
	InputScopes   []string
	OutputScopes  []string
	AdminOnly     bool
	StaffOnly     bool
}

func manageCommand(commandUpdate utils.CommandUpdate) {
	switch commandUpdate.Update.Message.Command() {
	case "credits":
		cmd := ShortCommands["credits"].ToCommand(commandUpdate, "The source code, available on GitHub at MoraGames/clockyuwu, is written entirely in GoLang and makes use of the \"telegram-bot-api\" library.\nFor any bug reports or feature proposals, please refer to the GitHub project.\n\nDeveloper:\n- Telegram: @MoraGames\n- Discord: @moragames\n- Instagram: @moragames.dev\n- GitHub: MoraGames\n\nProject:\n- Telegram: @clockygame_bot\n- GitHub: MoraGames/clockygame\n\nSpecial thanks go to the first testers (as well as players) of the minigame managed by the bot, \"Vano\", \"Ale\" and \"Alex\".")
		cmd.Execute()
	case "help":
		cmdsList := make([]string, 0)
		for _, scmd := range ShortCommands {
			cmdsList = append(cmdsList, scmd.ShortFormat())
		}

		cmd := ShortCommands["help"].ToCommand(commandUpdate, fmt.Sprintf("Bot Name: %v\nApp Name: %v\nApp Version: %v\n\nThis is a list of all possible commands within the bot:\n\n%v", App.BotAPI.Self.UserName, App.Name, App.Version, strings.Join(cmdsList, "\n")))
		cmd.Execute()
	case "ping":
		cmd := ShortCommands["ping"].ToCommand(commandUpdate, "pong")
		cmd.Execute()
	case "start":
		cmd := ShortCommands["start"].ToCommand(commandUpdate, fmt.Sprintf("%v is a bot that allows you to play a time-wasting game with one or more groups of friends within Telegram groups. Once the bot is added, the game mainly (but not exclusively) involves sending messages in the \"hh:mm\" format at certain times of the day, in exchange for valuable points. The person who has earned the most points at the end of the championship will be the new Clocky Champion!\nUse /help to get a list of all commands or /credits for more information about the project.\n\n- %v, an app from %v.", App.BotAPI.Self.UserName, App.Name, App.Author))
		cmd.Execute()
	default:
		//respond with an error message
		msg := tgbotapi.NewMessage(commandUpdate.Update.Message.Chat.ID, "Unknown command.\nUse /help to get a list of all commands.")
		msg.ReplyToMessageID = commandUpdate.Update.Message.MessageID
		App.BotAPI.Send(msg)

		App.Logger.WithFields(logrus.Fields{
			"message": commandUpdate.Update.Message.Text,
			"sender":  commandUpdate.Update.Message.From.UserName,
			"chat":    commandUpdate.Update.Message.Chat.Title,
		}).Debug("Response to unknown command sent successfully")
	}
}

func (cmd Command) Execute() error {
	//TODO: validate input scopes, output scopes, admin only and staff only properties

	//respond with the message provided by the command
	msg := tgbotapi.NewMessage(cmd.CommandUpdate.Update.Message.Chat.ID, cmd.Message)
	msg.ReplyToMessageID = cmd.CommandUpdate.Update.Message.MessageID
	App.BotAPI.Send(msg)

	App.Logger.WithFields(logrus.Fields{
		"message": cmd.CommandUpdate.Update.Message.Text,
		"sender":  cmd.CommandUpdate.Update.Message.From.UserName,
		"chat":    cmd.CommandUpdate.Update.Message.Chat.Title,
	}).Debugf("Response to \"%s\" command sent successfully", cmd.Command)

	return nil
}
