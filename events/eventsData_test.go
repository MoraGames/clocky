package events

import (
	"slices"
	"testing"
)

func TestRemoveValue(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		value string
		want  []string
	}{
		{"removes the only element", []string{"Add 2"}, "Add 2", []string{}},
		{"removes one of many", []string{"Add 2", "Sub 3", "Mul +2"}, "Sub 3", []string{"Add 2", "Mul +2"}},
		{"keeps everything when the value is absent", []string{"Add 2"}, "Sub 3", []string{"Add 2"}},
		{"removes every occurrence", []string{"Add 2", "Add 2"}, "Add 2", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveValue(tt.input, tt.value)
			if !slices.Equal(got, tt.want) {
				t.Errorf("RemoveValue(%v, %q) = %v, want %v", tt.input, tt.value, got, tt.want)
			}
			for _, v := range got {
				if v == "" {
					t.Errorf("RemoveValue returned an empty string: %v", got)
				}
			}
		})
	}
}
