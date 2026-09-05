package events

import (
	"testing"
	"time"

	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestButtonComboResolution(t *testing.T) {
	event, first := newButtonComboEvent(0)
	event.Activate(first, time.Now(), time.Now(), 2)

	winnerID, bonus, resolved := event.TryResolveButtonCombo(first.TelegramID)
	if !resolved || winnerID != first.TelegramID || bonus != 3 {
		t.Fatalf("expected first participant to win with +3, got winner=%d bonus=%d resolved=%t", winnerID, bonus, resolved)
	}
	if _, _, resolved = event.TryResolveButtonCombo(2); resolved {
		t.Fatal("expected a resolved button combo to reject later callbacks")
	}
}

func TestButtonComboMismatchAndExpiration(t *testing.T) {
	event, first := newButtonComboEvent(42)
	event.Activate(first, time.Now(), time.Now(), 2)

	winnerID, bonus, resolved := event.TryResolveButtonCombo(2)
	if !resolved || winnerID != first.TelegramID || bonus != -3 {
		t.Fatalf("expected first participant to win with -3, got winner=%d bonus=%d resolved=%t", winnerID, bonus, resolved)
	}

	event, _ = newButtonComboEvent(42)
	event.Activate(first, time.Now(), time.Now(), 2)
	messageID, expired := event.ExpireButtonCombo()
	if !expired || messageID != 42 {
		t.Fatalf("expected expiration to return message id, got id=%d expired=%t", messageID, expired)
	}
	if _, _, resolved = event.TryResolveButtonCombo(2); resolved {
		t.Fatal("expected an expired button combo to reject callbacks")
	}
}

func newButtonComboEvent(messageID int) (*Event, *structs.User) {
	telegramFirst := telegramUser(1)
	first := structs.NewUser(&telegramFirst)
	event := &Event{
		ButtonCombo:          true,
		ButtonComboMessageID: messageID,
		Partecipations:       make(map[int64]*EventPartecipation),
	}
	return event, first
}

func telegramUser(id int64) tgbotapi.User {
	return tgbotapi.User{ID: id, UserName: "user"}
}
