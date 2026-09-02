package events

func IsValidEventMessage(event *Event, text string) bool {
	if event.JokerFormat == "" {
		return text == event.Name
	}
	format, found := JokerFormatByName(event.JokerFormat)
	return found && format.Matches(text, event.Time)
}
