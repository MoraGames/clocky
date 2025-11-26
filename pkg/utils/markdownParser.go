package utils

import (
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MarkdownState struct {
	Type   string
	Offset int
}

const (
	MarkdownBoldDelimiter          = "**"
	MarkdownItalicDelimiter        = "%%"
	MarkdownUnderlineDelimiter     = "__"
	MarkdownStrikethroughDelimiter = "~~"
	MarkdownSpoilerDelimiter       = "||"
	MarkdownCodeDelimiter          = "```"

	MarkdownBoldEntityType          = "bold"
	MarkdownItalicEntityType        = "italic"
	MarkdownUnderlineEntityType     = "underline"
	MarkdownStrikethroughEntityType = "strikethrough"
	MarkdownSpoilerEntityType       = "spoiler"
	MarkdownCodeEntityType          = "code"
)

var (
	Markdowns = map[string]string{
		MarkdownBoldEntityType:          MarkdownBoldDelimiter,
		MarkdownItalicEntityType:        MarkdownItalicDelimiter,
		MarkdownUnderlineEntityType:     MarkdownUnderlineDelimiter,
		MarkdownStrikethroughEntityType: MarkdownStrikethroughDelimiter,
		MarkdownSpoilerEntityType:       MarkdownSpoilerDelimiter,
		MarkdownCodeEntityType:          MarkdownCodeDelimiter,
	}
)

func ParseToEntities(rawText string) ([]tgbotapi.MessageEntity, string) {
	var entities []tgbotapi.MessageEntity
	var states []MarkdownState
	var text []uint16

	utf16RawText := utf16.Encode([]rune(rawText))
	for i := 0; i < len(utf16RawText); {
		var delimiterFound string
		for entityType, delimiter := range Markdowns {
			if CheckFullDelimiter(utf16RawText, i, delimiter) {
				if len(states) > 0 && states[len(states)-1].Type == entityType { // closing tag
					entities = append(entities, tgbotapi.MessageEntity{
						Type:   entityType,
						Offset: states[len(states)-1].Offset,
						Length: len(text) - states[len(states)-1].Offset,
					})
					states = states[:len(states)-1]
				} else { // opening tag
					states = append(states, MarkdownState{
						Type:   entityType,
						Offset: len(text),
					})
				}
				delimiterFound = delimiter
				break
			}
		}

		if delimiterFound == "" {
			text = append(text, utf16RawText[i])
			i++
		} else {
			i += len(delimiterFound)
		}
	}

	return entities, string(utf16.Decode(text))
}

func CheckFullDelimiter(utf16Text []uint16, index int, delimiter string) bool {
	if index+len(delimiter)-1 >= len(utf16Text) {
		return false
	}
	for dc := 0; dc < len(delimiter); dc++ {
		if utf16Text[index+dc] != uint16(delimiter[dc]) {
			return false
		}
	}
	return true
}
