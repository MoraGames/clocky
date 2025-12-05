package utils

import (
	"regexp"
	"strconv"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	BasicMarkdown struct {
		Symbol string
		Type   string
	}
	BasicMarkdownState struct {
		Markdown          BasicMarkdown
		StartIndex        int
		BurnedIndexes     int
		DelimitersCounted int
	}

	OpenMarkdownLink struct {
		inTextUntil           int
		insideDelimitersCount int
		jumpAfterIndex        int
		entity                tgbotapi.MessageEntity
	}
)

const (
	MarkdownBoldDelimiter          = "**"
	MarkdownItalicDelimiter        = "%%"
	MarkdownUnderlineDelimiter     = "__"
	MarkdownStrikethroughDelimiter = "~~"
	MarkdownSpoilerDelimiter       = "||"
	MarkdownCodeDelimiter          = "`"
	MarkdownPreDelimiter           = "```"
	MarkdownTextLinkDelimiter      = "!"
	MarkdownTextMentionDelimiter   = "@"

	MarkdownBoldEntityType          = "bold"
	MarkdownItalicEntityType        = "italic"
	MarkdownUnderlineEntityType     = "underline"
	MarkdownStrikethroughEntityType = "strikethrough"
	MarkdownSpoilerEntityType       = "spoiler"
	MarkdownCodeEntityType          = "code"
	MarkdownPreEntityType           = "pre"
	MarkdownTextLinkEntityType      = "text_link"
	MarkdownTextMentionEntityType   = "text_mention"
)

var (
	basicMarkdowns = []BasicMarkdown{
		{Symbol: MarkdownBoldDelimiter, Type: MarkdownBoldEntityType},
		{Symbol: MarkdownItalicDelimiter, Type: MarkdownItalicEntityType},
		{Symbol: MarkdownUnderlineDelimiter, Type: MarkdownUnderlineEntityType},
		{Symbol: MarkdownStrikethroughDelimiter, Type: MarkdownStrikethroughEntityType},
		{Symbol: MarkdownSpoilerDelimiter, Type: MarkdownSpoilerEntityType},
	}
	basicMarkdownsStack = make([]BasicMarkdownState, 0)

	idRegexp = regexp.MustCompile(`^[0-9]+$`)
)

func ParseToEntities(rawText string, usersList []*tgbotapi.User) ([]tgbotapi.MessageEntity, string) {
	var entities []tgbotapi.MessageEntity
	var text []uint16

	utf16RawText := utf16.Encode([]rune(rawText))
	i := 0

	write := func(cp []uint16) {
		text = append(text, cp...)
		i += len(cp)
	}
	hasAheadIndex := func(str string, si int) bool {
		if si+len(str)-1 >= len(utf16RawText) {
			return false
		}
		for stri := 0; stri < len(str); stri++ {
			if utf16RawText[si+stri] != uint16(str[stri]) {
				return false
			}
		}
		return true
	}
	findAheadIndex := func(str string, si int) int {
		for ei := si; ei < len(utf16RawText); ei++ {
			if hasAheadIndex(str, ei) {
				return ei
			}
		}
		return -1
	}

	currentOpenLink := OpenMarkdownLink{
		inTextUntil:           -1,
		insideDelimitersCount: 0,
		jumpAfterIndex:        -1,
		entity:                tgbotapi.MessageEntity{},
	}
	markdownDelimiterCounter := 0

	for i = 0; i < len(utf16RawText); {
		if currentOpenLink.inTextUntil == -1 {
			//pre
			if hasAheadIndex(MarkdownPreDelimiter, i) {
				ei := findAheadIndex(MarkdownPreDelimiter, i+len(MarkdownPreDelimiter))
				if ei == -1 {
					//is not a pre entity
					//still can be a code entity (reason to shift by 2, not 3)
					write(utf16RawText[i : i+len(MarkdownPreDelimiter)-len(MarkdownCodeDelimiter)])
					continue
				}

				//extract language if any
				lang := ""
				j := i + len(MarkdownPreDelimiter)
				for ; j < ei; j++ {
					if utf16RawText[j] == uint16(' ') {
						lang = string(utf16.Decode(utf16RawText[i+len(MarkdownPreDelimiter) : j]))
						break
					}
				}

				//add opening markdown delimiter in the counter
				markdownDelimiterCounter += len(MarkdownPreDelimiter) + len(lang) + 1

				//define the pre entity
				entities = append(entities, tgbotapi.MessageEntity{
					Type:     "pre",
					Offset:   j + 1 - markdownDelimiterCounter,
					Length:   ei - (j + 1),
					Language: lang,
				})

				//add closing markdown delimiter in the counter and burn indexes for open basic markdowns
				markdownDelimiterCounter += len(MarkdownPreDelimiter)
				for bmdi := 0; bmdi < len(basicMarkdownsStack); bmdi++ {
					basicMarkdownsStack[bmdi].BurnedIndexes += 2 * len(MarkdownPreDelimiter)
				}

				//write the pre content and move forward
				write(utf16RawText[j+1 : ei])
				i = ei + len(MarkdownPreDelimiter)
				continue
			}

			//code
			if hasAheadIndex(MarkdownCodeDelimiter, i) {
				ei := findAheadIndex(MarkdownCodeDelimiter, i+len(MarkdownCodeDelimiter))
				if ei == -1 {
					//is not a code string
					write(utf16RawText[i : i+len(MarkdownCodeDelimiter)])
					continue
				}

				//add opening markdown delimiter in the counter
				markdownDelimiterCounter += 1

				//define the code entity
				entities = append(entities, tgbotapi.MessageEntity{
					Type:   "code",
					Offset: i + 1 - markdownDelimiterCounter,
					Length: ei - (i + 1),
				})

				//add closing markdown delimiter in the counter and burn indexes for open basic markdowns
				markdownDelimiterCounter += len(MarkdownCodeDelimiter)
				for bmdi := 0; bmdi < len(basicMarkdownsStack); bmdi++ {
					basicMarkdownsStack[bmdi].BurnedIndexes += 2 * len(MarkdownCodeDelimiter)
				}

				//write the code content and move forward
				write(utf16RawText[i+len(MarkdownCodeDelimiter) : ei])
				i = ei + len(MarkdownCodeDelimiter)
				continue
			}

			//text link or text mention
			if hasAheadIndex("[", i) {
				mi := findAheadIndex("](", i+1)
				ei := findAheadIndex(")", mi+2)
				if mi == -1 || ei == -1 {
					//is not a text link or text mention
					write(utf16RawText[i : i+1])
					continue
				}

				//add opening markdown delimiter in the counter
				markdownDelimiterCounter += 1

				//initialize the open link state
				currentOpenLink.inTextUntil = mi
				currentOpenLink.jumpAfterIndex = ei + 1
				currentOpenLink.insideDelimitersCount = 0

				//check if is text link or text mention
				if hasAheadIndex(MarkdownTextMentionDelimiter, mi+2) {
					//text mention

					//retrieve the text mention tag
					tag := string(utf16.Decode(utf16RawText[mi+2+len(MarkdownTextMentionDelimiter) : ei]))
					isUserId := idRegexp.MatchString(tag)

					var userId int64
					var username string
					if isUserId {
						//parse user's id
						var err error
						userId, err = strconv.ParseInt(tag, 10, 64)
						if err != nil {
							panic(err)
						}
					} else {
						//parse user's username
						username = tag
					}

					//search for the user in the known users
					var mentionedUser *tgbotapi.User
					for _, u := range usersList {
						if u == nil {
							continue
						}
						if (isUserId && u.ID == userId) || (!isUserId && u.UserName != "" && u.UserName == username) {
							mentionedUser = u
							break
						}
					}

					//valorize the mention entity
					currentOpenLink.entity = tgbotapi.MessageEntity{
						Type:   MarkdownTextMentionEntityType,
						Offset: i + 1 - markdownDelimiterCounter,
						Length: mi - (i + 1),
						User:   mentionedUser,
					}
				} else if hasAheadIndex(MarkdownTextLinkDelimiter, mi+2) {
					//text link

					//valorize the link entity
					currentOpenLink.entity = tgbotapi.MessageEntity{
						Type:   MarkdownTextLinkEntityType,
						Offset: i + 1 - markdownDelimiterCounter,
						Length: mi - (i + 1),
						URL:    string(utf16.Decode(utf16RawText[mi+2+len(MarkdownTextLinkDelimiter) : ei])),
					}
				}

				//move forward to the text of the link/mention
				i++
				continue
			}
		}

		//end of text link or text mention
		if currentOpenLink.inTextUntil != -1 && i == currentOpenLink.inTextUntil {
			//add closing markdown delimiter in the counter
			markdownDelimiterCounter += 3 + ((currentOpenLink.jumpAfterIndex - 1) - (currentOpenLink.inTextUntil + 2))

			//update the current open link state
			currentOpenLink.entity.Length -= currentOpenLink.insideDelimitersCount
			entities = append(entities, currentOpenLink.entity)

			//move forward after the link/mention
			i = currentOpenLink.jumpAfterIndex

			//reset the current open link state
			currentOpenLink = OpenMarkdownLink{
				inTextUntil:           -1,
				insideDelimitersCount: 0,
				jumpAfterIndex:        -1,
				entity:                tgbotapi.MessageEntity{},
			}
			continue
		}

		//basic markdowns
		var neededContinue bool
		for _, bmd := range basicMarkdowns {
			if hasAheadIndex(bmd.Symbol, i) {
				ei := findAheadIndex(bmd.Symbol, i+len(bmd.Symbol))
				if ei == -1 && !(len(basicMarkdownsStack) > 0 && basicMarkdownsStack[len(basicMarkdownsStack)-1].Markdown.Type == bmd.Type) {
					// Isn't this markdown
					write(utf16RawText[i : i+len(bmd.Symbol)])
					neededContinue = true
					break
				}

				//check if is opening or closing markdown
				if len(basicMarkdownsStack) > 0 && basicMarkdownsStack[len(basicMarkdownsStack)-1].Markdown.Type == bmd.Type {
					//closing markdown

					//retrieve and remove last basic markdown status from the stack and create the entity
					bmds := basicMarkdownsStack[len(basicMarkdownsStack)-1]
					basicMarkdownsStack = basicMarkdownsStack[:len(basicMarkdownsStack)-1]
					entities = append(entities, tgbotapi.MessageEntity{
						Type:   bmds.Markdown.Type,
						Offset: bmds.StartIndex - bmds.DelimitersCounted,
						Length: i - bmds.StartIndex - bmds.BurnedIndexes,
					})

					//add closing markdown delimiter in the counter and burn indexes for open basic markdowns
					for bmdi := 0; bmdi < len(basicMarkdownsStack); bmdi++ {
						basicMarkdownsStack[bmdi].BurnedIndexes += len(bmd.Symbol)
					}
					currentOpenLink.insideDelimitersCount += len(bmd.Symbol)
					markdownDelimiterCounter += len(bmd.Symbol)
				} else {
					//opening markdown

					//add opening markdown delimiter in the counter and burn indexes for open basic markdowns
					markdownDelimiterCounter += len(bmd.Symbol)
					currentOpenLink.insideDelimitersCount += len(bmd.Symbol)
					for bmdi := 0; bmdi < len(basicMarkdownsStack); bmdi++ {
						basicMarkdownsStack[bmdi].BurnedIndexes += len(bmd.Symbol)
					}

					//push to the stack the basic markdown status
					basicMarkdownsStack = append(basicMarkdownsStack, BasicMarkdownState{
						Markdown:          bmd,
						StartIndex:        i + len(bmd.Symbol),
						DelimitersCounted: markdownDelimiterCounter,
						BurnedIndexes:     0,
					})
				}

				//move forward
				i += len(bmd.Symbol)
				neededContinue = true
				break
			}
		}
		if neededContinue {
			continue
		}

		//write normal character
		write(utf16RawText[i : i+1])
	}

	return entities, string(utf16.Decode(text))
}
