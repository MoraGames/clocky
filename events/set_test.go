package events

import (
	"testing"
	"time"
)

func TestCurrentDailyExpiration_BeforeCutoffUsesToday(t *testing.T) {
	loc := time.FixedZone("CET", 3600)
	now := time.Date(2026, 8, 1, 23, 59, 29, 0, loc)
	got := currentDailyExpiration(now)
	want := time.Date(2026, 8, 1, 23, 59, 30, 0, loc)

	if !got.Equal(want) {
		t.Fatalf("expected expiration %v, got %v", want, got)
	}
}

func TestCurrentDailyExpiration_AfterCutoffUsesTomorrow(t *testing.T) {
	loc := time.FixedZone("CET", 3600)
	now := time.Date(2026, 8, 1, 23, 59, 31, 0, loc)
	got := currentDailyExpiration(now)
	want := time.Date(2026, 8, 2, 23, 59, 30, 0, loc)

	if !got.Equal(want) {
		t.Fatalf("expected expiration %v, got %v", want, got)
	}
}

func TestMergeSetsWithDefaults_PreservesMissingDefaultSets(t *testing.T) {
	persisted := SetJsonSlice{
		{Name: "Short Equal", Pattern: "?a:aa", Typology: "static", Enabled: true},
		{Name: "Repeat", Pattern: "ab:ab", Typology: "static", Enabled: false},
	}

	merged := mergeSetsWithDefaults(persisted)
	if len(merged) != len(defaultSets()) {
		t.Fatalf("expected %d merged sets, got %d", len(defaultSets()), len(merged))
	}

	var foundHalf bool
	for _, set := range merged {
		if set.Name == "Half" {
			foundHalf = true
			if set.Verify == nil {
				t.Fatal("expected Half to keep its verifier")
			}
		}
	}

	if !foundHalf {
		t.Fatal("expected merged sets to include Half")
	}
}
