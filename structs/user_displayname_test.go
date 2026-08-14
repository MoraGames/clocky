package structs

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestDisplayNameFallsBackToName(t *testing.T) {
	u := &tgbotapi.User{FirstName: "Chiaraa", UserName: ""}
	if got := DisplayName(u); got != "Chiaraa" {
		t.Errorf("DisplayName con solo il nome = %q, atteso %q", got, "Chiaraa")
	}
}

func TestDisplayNamePrefersUsername(t *testing.T) {
	u := &tgbotapi.User{FirstName: "Chiaraa", UserName: "chiara99"}
	if got := DisplayName(u); got != "chiara99" {
		t.Errorf("DisplayName con @username = %q, atteso %q", got, "chiara99")
	}
}

// NewUser deve dare il nome anche a chi non ha @username: è la creazione,
// gia' corretta in precedenza. Il regression vero e' sul ramo di update.
func TestNewUserGetsNameWithoutUsername(t *testing.T) {
	u := NewUser(&tgbotapi.User{ID: 1, FirstName: "Chiaraa"})
	if u.UserName != "Chiaraa" {
		t.Errorf("NewUser senza @username = %q, atteso %q", u.UserName, "Chiaraa")
	}
}
