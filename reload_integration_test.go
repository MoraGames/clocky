package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/structs"
	"github.com/sirupsen/logrus"
)

func setupTempGameFiles(t *testing.T) types.Utils {
	t.Helper()

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

	root, err := os.OpenRoot("files")
	if err != nil {
		t.Fatalf("open files root: %v", err)
	}

	oldRoot := App.FilesRoot
	App.FilesRoot = root
	t.Cleanup(func() {
		if err := root.Close(); err != nil {
			t.Fatalf("close files root: %v", err)
		}
		App.FilesRoot = oldRoot
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	return types.Utils{Logger: logger, TimeFormat: time.RFC3339}
}

func writeTempJSON(t *testing.T, fileName string, value any) {
	t.Helper()

	file, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		t.Fatalf("marshal %s: %v", fileName, err)
	}

	if err := os.WriteFile("files/"+fileName, file, 0o644); err != nil {
		t.Fatalf("write %s: %v", fileName, err)
	}
}

func hasSet(setSlice events.SetSlice, setName string) bool {
	for _, set := range setSlice {
		if set != nil && set.Name == setName {
			return true
		}
	}
	return false
}

func TestReloadStatus_ExpiredFilesFallBackToDefaults(t *testing.T) {
	utils := setupTempGameFiles(t)

	events.Sets = nil
	events.SetsJson = events.SetFile{}
	events.CurrentChampionship = nil

	writeTempJSON(t, "sets.json", events.SetFile{
		Slice: events.SetJsonSlice{
			{Name: "Short Equal", Pattern: "?a:aa", Typology: "static", Enabled: true},
		},
		Expiration: time.Now().Add(-time.Hour),
	})
	writeTempJSON(t, "championship.json", &structs.Championship{
		Name:         "Clocky Championship",
		StartDate:    time.Now().Add(-14 * 24 * time.Hour),
		Duration:     14 * 24 * time.Hour,
		Expiration:   time.Now().Add(-time.Hour),
		Status:       "ended",
		FinalRanking: nil,
	})

	reloadStatus([]types.Reload{
		{
			FileName:   "sets.json",
			DataStruct: &events.SetsJson,
			Validate:   events.SetsFileValid,
			IfOkay:     events.AssignSetsFromSetsJson,
			IfFail:     events.AssignSetsWithDefault,
		},
		{
			FileName:   "championship.json",
			DataStruct: &events.CurrentChampionship,
			Validate:   events.ChampionshipFileValid,
			IfOkay:     events.NormalizeChampionshipData,
			IfFail:     events.AssignChampionshipWithDefault,
		},
	}, utils)

	if len(events.Sets) != len(events.SetsFunctions) {
		t.Fatalf("expected default set count %d, got %d", len(events.SetsFunctions), len(events.Sets))
	}
	if !hasSet(events.Sets, "Half") {
		t.Fatal("expected fallback sets to include Half")
	}
	if events.CurrentChampionship == nil {
		t.Fatal("expected championship fallback to initialize data")
	}
	if events.CurrentChampionship.IsExpired(time.Now()) {
		t.Fatal("expected championship fallback to be valid")
	}
}

func TestEventsResetWritesFilesAndReloadStatusLoadsThem(t *testing.T) {
	utils := setupTempGameFiles(t)

	events.AssignSetsWithDefault(utils)
	events.Events = events.NewEventsData(false, utils)
	events.Events.Reset(false, nil, utils)

	var savedSets events.SetFile
	if file, err := os.ReadFile("files/sets.json"); err != nil {
		t.Fatalf("read sets file: %v", err)
	} else if err := json.Unmarshal(file, &savedSets); err != nil {
		t.Fatalf("unmarshal sets file: %v", err)
	}
	if savedSets.Expiration.IsZero() {
		t.Fatal("expected sets file to persist expiration")
	}
	if time.Now().After(savedSets.Expiration) {
		t.Fatal("expected sets file expiration to be in the future")
	}

	var savedEvents events.EventsData
	if file, err := os.ReadFile("files/events.json"); err != nil {
		t.Fatalf("read events file: %v", err)
	} else if err := json.Unmarshal(file, &savedEvents); err != nil {
		t.Fatalf("unmarshal events file: %v", err)
	}
	if savedEvents.Expiration.IsZero() {
		t.Fatal("expected events file to persist expiration")
	}
	if time.Now().After(savedEvents.Expiration) {
		t.Fatal("expected events file expiration to be in the future")
	}

	events.Sets = nil
	events.SetsJson = events.SetFile{}
	events.Events = nil

	reloadStatus([]types.Reload{
		{
			FileName:   "sets.json",
			DataStruct: &events.SetsJson,
			Validate:   events.SetsFileValid,
			IfOkay:     events.AssignSetsFromSetsJson,
			IfFail:     events.AssignSetsWithDefault,
		},
		{
			FileName:   "events.json",
			DataStruct: &events.Events,
			Validate:   events.EventsFileValid,
			IfOkay:     events.NormalizeEventsData,
			IfFail:     events.AssignEventsWithDefault,
		},
	}, utils)

	if len(events.Sets) != len(events.SetsFunctions) {
		t.Fatalf("expected reloaded set count %d, got %d", len(events.SetsFunctions), len(events.Sets))
	}
	if !hasSet(events.Sets, "Half") {
		t.Fatal("expected reloaded sets to include Half")
	}
	if events.Events == nil {
		t.Fatal("expected events to reload from saved file")
	}
	if events.Events.Expiration.IsZero() {
		t.Fatal("expected reloaded events to carry expiration")
	}
	if time.Now().After(events.Events.Expiration) {
		t.Fatal("expected reloaded events expiration to be in the future")
	}
}
