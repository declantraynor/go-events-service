package usecases

import (
	"fmt"
)

type InvalidTimestampError struct {
	timestamp  string
	notISO8601 bool
	notUTC     bool
}

func (err InvalidTimestampError) Error() string {
	var format string
	if err.notISO8601 {
		format = "timestamp:%q does not conform to ISO8601"
	} else if err.notUTC {
		format = "timestamp:%q is not UTC"
	}
	return fmt.Sprintf(format, err.timestamp)
}

type InvalidTimeRangeError struct {
	from string
	to   string
}

func (err InvalidTimeRangeError) Error() string {
	return fmt.Sprintf("from:%q is later than to:%q", err.from, err.to)
}
