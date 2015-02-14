// Package usecases implements application-specific business logic
// for the events service.
package usecases

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/declantraynor/go-events-service/domain"
)

// SanitizeName returns its argument with all white space removed.
// Leading and trailing white space is trimmed. Any intervening
// white space characters are replaced with a dash.
func SanitizeName(name string) string {
	re := regexp.MustCompile("\\s+")
	return re.ReplaceAllString(strings.TrimSpace(name), "-")
}

// ParseTimestamp attempts to parse an ISO8601 (RFC3339) compliant
// UTC time value from its argument string. It returns a time.Time
// value as well as any error encountered.
func ParseTimestamp(timestamp string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, errors.New("timestamps must conform to ISO8601")
	}
	if _, utcOffset := t.Zone(); utcOffset != 0 {
		return time.Time{}, errors.New("timestamps must be UTC")
	}
	return t, nil
}

type EventInteractor struct {
	repo domain.EventRepo
}

func (interactor *EventInteractor) Add(name string, timestamp string) (domain.Event, error) {
	if _, err := ParseTimestamp(timestamp); err != nil {
		return domain.Event{}, err
	}
	event := domain.Event{Name: SanitizeName(name), Timestamp: timestamp}
	return interactor.repo.Store(event)
}
