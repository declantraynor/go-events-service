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
		format = "%s does not conform to ISO8601"
	} else if err.NotUTC {
		format = "%s is not UTC"
	}
	return fmt.Sprintf(format, err.Timestamp)
}

type InvalidTimeRangeError struct {
	From string
	To   string
}

func (err InvalidTimeRangeError) Error() string {
	return fmt.Sprintf("%s is later than %s", err.From, err.To)
}
