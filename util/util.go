package util

import (
	"fmt"
	"time"

	"willnorris.com/go/newbase60"
)

const day = 24 * time.Hour

var epoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// EpochDaysToTime converts the provided newbase60-encoded epoch days into a
// Time.  Epoch days are the number of days that have elapsed since 1970-01-01.
func EpochDaysToTime(s string) time.Time {
	n := newbase60.DecodeToInt(s)
	return epoch.Add(time.Duration(n) * day)
}

// TimeToEpochDays converts t into epoch days, encoded in newbase60.
func TimeToEpochDays(t time.Time) string {
	d := t.Sub(epoch) / day
	s := newbase60.EncodeInt(int(d))
	return fmt.Sprintf("%03s", s)
}
