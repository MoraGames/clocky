package events

import (
	"math/rand"
	"strings"
	"time"
)

type JokerFormat struct {
	Name        string
	Description string
	Format      func(time.Time) string
	Parse       func(string) (int, int, bool)
}

var jokerFormats = []JokerFormat{
	{Name: "Emoji Time", Description: "Usa le emoji con i numeri da 0️⃣ a 9️⃣ per formare l'orario!", Format: formatKeycap, Parse: parseKeycap},
	{Name: "Roman Numerals", Description: "Ricorda, solo N, I, V e X erano simboli utilizzati; dovrai separarli con spazi.", Format: formatRoman, Parse: parseRoman},
	{Name: "Many Repetition", Description: "Un punto interrogativo, ripetuto per molte domande... oppure diversi orari e qualche #.", Format: formatRepeatedQuestionMarks, Parse: parseRepeatedQuestionMarks},
	{Name: "Clocky Clock", Description: "Semplicemente 4 orologi, tutti che segnano l'ora in punto!", Format: formatClockyClock, Parse: parseClockyClock},
	{Name: "Alternative Sings", Description: "Un layout, 10 cifre alternative, tutte che vanno dal ! al =.", Format: formatAlternativeSings, Parse: parseAlternativeSings},
}

func JokerFormats() []JokerFormat {
	return append([]JokerFormat(nil), jokerFormats...)
}

func JokerFormatByName(name string) (JokerFormat, bool) {
	for _, format := range jokerFormats {
		if format.Name == name {
			return format, true
		}
	}
	return JokerFormat{}, false
}

func (ed *EventsData) AssignJokerFormats() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, event := range ed.Map {
		event.JokerFormat = ""
		if !event.Enabled {
			continue
		}
		for _, set := range Sets {
			if set.Enabled && set.Joker && set.Verify != nil && set.Verify(SplitTime(event.Time)) {
				event.JokerFormat = jokerFormats[random.Intn(len(jokerFormats))].Name
				break
			}
		}
	}
}

func (format JokerFormat) Matches(text string, expected time.Time) bool {
	hour, minute, ok := format.Parse(strings.TrimSpace(text))
	return ok && hour == expected.Hour() && minute == expected.Minute()
}

var keycapDigits = []string{"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}

func formatKeycap(value time.Time) string {
	digits := []int{value.Hour() / 10, value.Hour() % 10, value.Minute() / 10, value.Minute() % 10}
	parts := make([]string, 0, len(digits))
	for _, digit := range digits {
		parts = append(parts, keycapDigits[digit])
	}
	return parts[0] + parts[1] + " : " + parts[2] + parts[3]
}

func parseKeycap(text string) (int, int, bool) {
	parts := strings.Split(text, " : ")
	if len(parts) != 2 {
		return 0, 0, false
	}
	digits := make([]int, 0, 4)
	for _, part := range parts {
		for len(part) > 0 {
			found := false
			for digit, symbol := range keycapDigits {
				if strings.HasPrefix(part, symbol) {
					digits = append(digits, digit)
					part = strings.TrimPrefix(part, symbol)
					found = true
					break
				}
			}
			if !found {
				return 0, 0, false
			}
		}
	}
	if len(digits) != 4 {
		return 0, 0, false
	}
	hour, minute := digits[0]*10+digits[1], digits[2]*10+digits[3]
	return hour, minute, hour < 24 && minute < 60
}

var romanDigits = []string{"N", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}

func formatRoman(value time.Time) string {
	digits := []int{value.Hour() / 10, value.Hour() % 10, value.Minute() / 10, value.Minute() % 10}
	parts := make([]string, 0, len(digits))
	for _, digit := range digits {
		parts = append(parts, romanDigits[digit])
	}
	return parts[0] + " " + parts[1] + " : " + parts[2] + " " + parts[3]
}

func parseRoman(text string) (int, int, bool) {
	parts := strings.Split(text, " : ")
	if len(parts) != 2 {
		return 0, 0, false
	}
	digits := make([]int, 0, 4)
	for _, part := range parts {
		for _, token := range strings.Fields(part) {
			found := false
			for digit, symbol := range romanDigits {
				if token == symbol {
					digits = append(digits, digit)
					found = true
					break
				}
			}
			if !found {
				return 0, 0, false
			}
		}
	}
	if len(digits) != 4 {
		return 0, 0, false
	}
	hour, minute := digits[0]*10+digits[1], digits[2]*10+digits[3]
	return hour, minute, hour < 24 && minute < 60
}

func formatRepeatedQuestionMarks(value time.Time) string {
	digits := []int{value.Hour() / 10, value.Hour() % 10, value.Minute() / 10, value.Minute() % 10}
	parts := make([]string, 0, len(digits))
	for _, digit := range digits {
		parts = append(parts, repeatedQuestionMarks(digit))
	}
	return parts[0] + " " + parts[1] + " : " + parts[2] + " " + parts[3]
}

func repeatedQuestionMarks(digit int) string {
	if digit == 0 {
		return "#"
	}
	return strings.Repeat("?", digit)
}

func parseRepeatedQuestionMarks(text string) (int, int, bool) {
	parts := strings.Split(text, " : ")
	if len(parts) != 2 {
		return 0, 0, false
	}
	digits := make([]int, 0, 4)
	for _, part := range parts {
		for _, token := range strings.Fields(part) {
			digit := 0
			if token == "#" {
				digits = append(digits, 0)
				continue
			}
			if len(token) > 9 || strings.Trim(token, "?") != "" {
				return 0, 0, false
			}
			digit = len(token)
			digits = append(digits, digit)
		}
	}
	if len(digits) != 4 {
		return 0, 0, false
	}
	hour, minute := digits[0]*10+digits[1], digits[2]*10+digits[3]
	return hour, minute, hour < 24 && minute < 60
}

var clockyClockDigits = []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘"}

func formatClockyClock(value time.Time) string {
	digits := []int{value.Hour() / 10, value.Hour() % 10, value.Minute() / 10, value.Minute() % 10}
	parts := make([]string, 0, len(digits))
	for _, digit := range digits {
		parts = append(parts, clockyClockDigits[digit])
	}
	return parts[0] + parts[1] + ":" + parts[2] + parts[3]
}

func parseClockyClock(text string) (int, int, bool) {
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	digits := make([]int, 0, 4)
	for _, part := range parts {
		for _, token := range strings.Fields(part) {
			found := false
			for digit, symbol := range clockyClockDigits {
				if token == symbol {
					digits = append(digits, digit)
					found = true
					break
				}
			}
			if !found {
				return 0, 0, false
			}
		}
	}
	if len(digits) != 4 {
		return 0, 0, false
	}
	hour, minute := digits[0]*10+digits[1], digits[2]*10+digits[3]
	return hour, minute, hour < 24 && minute < 60
}

var alternativeSingsDigits = []string{"!", "\"", "£", "$", "%", "&", "/", "(", ")", "="}

func formatAlternativeSings(value time.Time) string {
	digits := []int{value.Hour() / 10, value.Hour() % 10, value.Minute() / 10, value.Minute() % 10}
	parts := make([]string, 0, len(digits))
	for _, digit := range digits {
		parts = append(parts, alternativeSingsDigits[digit])
	}
	return parts[0] + parts[1] + ":" + parts[2] + parts[3]
}

func parseAlternativeSings(text string) (int, int, bool) {
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	digits := make([]int, 0, 4)
	for _, part := range parts {
		for _, token := range strings.Fields(part) {
			found := false
			for digit, symbol := range alternativeSingsDigits {
				if token == symbol {
					digits = append(digits, digit)
					found = true
					break
				}
			}
			if !found {
				return 0, 0, false
			}
		}
	}
	if len(digits) != 4 {
		return 0, 0, false
	}
	hour, minute := digits[0]*10+digits[1], digits[2]*10+digits[3]
	return hour, minute, hour < 24 && minute < 60
}
