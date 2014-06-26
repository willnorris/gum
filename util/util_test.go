package util

import (
	"testing"
	"time"
)

func TestEpochDayToTime(t *testing.T) {
	tests := []struct {
		s string
		t time.Time
	}{
		// time at epoch
		{"000", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
		// the largest date Go can handle
		{"VeB", time.Date(2262, 4, 11, 0, 0, 0, 0, time.UTC)},
		// today (the date this test was written)
		{"4Wn", time.Date(2014, 6, 26, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		if got, want := EpochDaysToTime(tt.s), tt.t; !got.Equal(want) {
			t.Errorf("EpochDaysToTime(%q) got: %v, want: %v", tt.s, got, want)
		}
		if got, want := TimeToEpochDays(tt.t), tt.s; got != want {
			t.Errorf("TimeToEpochDays(%q) got: %v, want: %v", tt.t, got, want)
		}
	}
}
