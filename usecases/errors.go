package usecases

import (
	"fmt"
)

type InvalidTimestampError struct {
	Timestamp  string
	NotISO8601 bool
	NotUTC     bool
}

func (err InvalidTimestampError) Error() string {
	var format string
	if err.NotISO8601 {
		format = "timestamp:%q does not conform to ISO8601"
	} else if err.NotUTC {
		format = "timestamp:%q is not UTC"
	}
	return fmt.Sprintf(format, err.Timestamp)
}

type InvalidTimeRangeError struct {
	From string
	To   string
}

func (err InvalidTimeRangeError) Error() string {
	return fmt.Sprintf("from:%q is later than to:%q", err.From, err.To)
}
