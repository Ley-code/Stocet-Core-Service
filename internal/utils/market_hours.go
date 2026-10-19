package utils

import "time"

// Ethiopian Securities Exchange trading window (Africa/Addis_Ababa, UTC+3, no DST).
// Open 09:30 EAT (inclusive), close 15:00 EAT (exclusive).
const (
	ESXMarketOpenMinutes  = 9*60 + 30 // 09:30
	ESXMarketCloseMinutes = 15 * 60   // 15:00
)

// IsEATMarketOpen reports whether `now` falls inside the ESX trading window.
// The function operates by converting `now` into Africa/Addis_Ababa local time
// and comparing the wall-clock minute-of-day against the open/close window.
//
// Falls back to a fixed UTC+3 offset if the IANA tz data isn't installed
// in the runtime image (matters on minimal container images).
func IsEATMarketOpen(now time.Time) bool {
	loc, err := time.LoadLocation("Africa/Addis_Ababa")
	if err != nil {
		loc = time.FixedZone("EAT", 3*60*60)
	}
	t := now.In(loc)
	total := t.Hour()*60 + t.Minute()
	return total >= ESXMarketOpenMinutes && total < ESXMarketCloseMinutes
}
