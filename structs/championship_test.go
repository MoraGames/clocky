package structs

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/sirupsen/logrus"
)

func TestChampionshipRefreshExpirationAndExpiry(t *testing.T) {
	loc := time.FixedZone("CET", 3600)
	start := time.Date(2026, 8, 1, 23, 59, 25, 0, loc)
	championship := CreateChampionship("Test Championship", start, 14*24*time.Hour)

	wantExpiration := start.Add(14 * 24 * time.Hour)
	if !championship.Expiration.Equal(wantExpiration) {
		t.Fatalf("expected expiration %v, got %v", wantExpiration, championship.Expiration)
	}

	if championship.IsExpired(wantExpiration.Add(-time.Second)) {
		t.Fatal("championship should not be expired just before expiration")
	}

	if !championship.IsExpired(wantExpiration) {
		t.Fatal("championship should be expired at expiration time")
	}

	championship.Expiration = time.Time{}
	championship.RefreshExpiration()
	if !championship.Expiration.Equal(wantExpiration) {
		t.Fatalf("expected refreshed expiration %v, got %v", wantExpiration, championship.Expiration)
	}
}

func TestChampionshipSaveAndReadRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	if err := os.MkdirAll("files", 0o755); err != nil {
		t.Fatalf("create files dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	utils := types.Utils{Logger: logger, TimeFormat: time.RFC3339}

	start := time.Date(2026, 8, 1, 23, 59, 25, 0, time.FixedZone("CET", 3600))
	championship := CreateChampionship("Clocky Championship", start, 14*24*time.Hour)
	championship.End([]Rank{{
		UserTelegramID: 123,
		Username:       "TestUser",
		Points:         99,
		Partecipations: 7,
	}})
	championship.SaveOnFile(utils)

	loaded := ReadFromFile(utils)
	if loaded == nil {
		t.Fatal("expected championship to be read back from file")
	}
	if loaded.Name != championship.Name {
		t.Fatalf("expected championship name %q, got %q", championship.Name, loaded.Name)
	}
	if !loaded.Expiration.Equal(championship.Expiration) {
		t.Fatalf("expected expiration %v, got %v", championship.Expiration, loaded.Expiration)
	}
	if loaded.Status != "ended" {
		t.Fatalf("expected ended championship, got %q", loaded.Status)
	}
	if len(loaded.FinalRanking) != 1 || loaded.FinalRanking[0].Username != "TestUser" {
		t.Fatal("expected final ranking to round-trip through file")
	}
}
