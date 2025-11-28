package types

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// IsCommand returns true if message starts with a "bot_command" entity in either Text or Caption
func IsCommand(m *tgbotapi.Message) bool {
	if m.Text != "" && m.IsCommand() {
		return true
	}

	if m.Caption != "" && IsCaptionCommand(m) {
		return true
	}

	return false
}

// IsCaptionCommand checks if caption starts with a bot_command entity
func IsCaptionCommand(m *tgbotapi.Message) bool {
	if len(m.CaptionEntities) == 0 {
		return false
	}

	entity := m.CaptionEntities[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

// Command returns the command from either Text or Caption
func Command(m *tgbotapi.Message) string {
	command := CommandWithAt(m)

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAtreturns the command from either Text or Caption with @ mention
func CommandWithAt(m *tgbotapi.Message) string {
	if !IsCommand(m) {
		return ""
	}

	if m.Text != "" && m.IsCommand() {
		entity := m.Entities[0]
		return m.Text[1:entity.Length]
	}

	if m.Caption != "" && IsCaptionCommand(m) {
		entity := m.CaptionEntities[0]
		return m.Caption[1:entity.Length]
	}

	return ""
}

// CommandArguments returns arguments from either Text or Caption
func CommandArguments(m *tgbotapi.Message) string {
	if !IsCommand(m) {
		return ""
	}

	if m.Text != "" && m.IsCommand() {
		entity := m.Entities[0]
		if len(m.Text) == entity.Length {
			return ""
		}
		return m.Text[entity.Length+1:]
	}

	if m.Caption != "" && IsCaptionCommand(m) {
		entity := m.CaptionEntities[0]
		if len(m.Caption) == entity.Length {
			return ""
		}
		return m.Caption[entity.Length+1:]
	}

	return ""
}
