// Package usecases implements application-specific business logic
// for the events service.
package usecases

import (
	"github.com/declantraynor/go-events-service/domain"
)

type EventInteractor struct {
	Store domain.EventStore
}

func (interactor *EventInteractor) Add(name string, timestamp string) error {

	parsedTimestamp, err := ParseTimestamp(timestamp)
	if err != nil {
		return err
	}

	event := domain.Event{Name: name, Timestamp: parsedTimestamp.Unix()}
	if err := interactor.Store.Put(event); err != nil {
		return err
	}

	return nil
}
