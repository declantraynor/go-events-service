package usecases

import (
	"time"
)

// ParseTimestamp attempts to parse an ISO8601 (RFC3339) compliant
// UTC time value from its argument string. It returns a time.Time
// value as well as any error encountered.
func ParseTimestamp(timestamp string) (time.Time, error) {

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, InvalidTimestampError{Timestamp: timestamp, NotISO8601: true}
	}

	if _, utcOffset := t.Zone(); utcOffset != 0 {
		return time.Time{}, InvalidTimestampError{Timestamp: timestamp, NotUTC: true}
	}

	return t, nil
}
