package events

import (
	"strings"
	"testing"
	"time"
)

func TestJokerFormatsRoundTrip(t *testing.T) {
	for _, format := range JokerFormats() {
		for hourValue := 0; hourValue < 24; hourValue++ {
			for minuteValue := 0; minuteValue < 60; minuteValue++ {
				expected := time.Date(0, time.January, 1, hourValue, minuteValue, 0, 0, time.UTC)
				text := format.Format(expected)
				parts := strings.Fields(text)
				if len(parts) != 5 || parts[2] != ":" {
					t.Fatalf("format %q did not use the canonical spaced format: %q", format.Name, text)
				}
				hour, minute, ok := format.Parse(text)
				if !ok || hour != expected.Hour() || minute != expected.Minute() {
					t.Fatalf("format %q failed round trip with %q", format.Name, text)
				}
				if !format.Matches(text, expected) {
					t.Fatalf("format %q did not match its own output %q", format.Name, text)
				}
			}
		}
	}
}

func TestJokerFormatsRejectCompactFormat(t *testing.T) {
	expected := time.Date(0, time.January, 1, 12, 34, 0, 0, time.UTC)
	for _, format := range JokerFormats() {
		text := format.Format(expected)
		compact := strings.ReplaceAll(text, " ", "")
		if format.Matches(compact, expected) {
			t.Fatalf("format %q accepted compact joker format %q", format.Name, compact)
		}
	}
}

func TestJokerFormatsRejectWrongTime(t *testing.T) {
	expected := time.Date(0, time.January, 1, 12, 34, 0, 0, time.UTC)
	for _, format := range JokerFormats() {
		if format.Matches(format.Format(expected), expected.Add(time.Minute)) {
			t.Fatalf("format %q accepted the wrong time", format.Name)
		}
	}
}

func TestRepeatedQuestionMarksUsesQuestionMarks(t *testing.T) {
	expected := time.Date(0, time.January, 1, 12, 34, 0, 0, time.UTC)
	format, found := JokerFormatByName("Many Repetition")
	if !found {
		t.Fatal("expected repeated question marks format")
	}
	if got := format.Format(expected); got != "? ?? : ??? ????" {
		t.Fatalf("unexpected repeated question marks format: %q", got)
	}
	if !format.Matches("? ?? : ??? ????", expected) {
		t.Fatal("expected question-mark format to match")
	}
}

func TestSetJsonPreservesJoker(t *testing.T) {
	set := &Set{Name: "Rise", Pattern: "ab:cd", Typology: "static", Enabled: true, Joker: true}
	jsonSet := SetSlice{set}.ToJsonSlice()[0]
	if !jsonSet.Joker {
		t.Fatal("expected joker flag in JSON representation")
	}

	restored := SetJsonSlice{jsonSet}.ToSlice()[0]
	if !restored.Joker {
		t.Fatal("expected joker flag after restoring JSON representation")
	}
}
