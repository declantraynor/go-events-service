package usecases

import (
	"errors"

	"github.com/declantraynor/go-events-service/domain"
)

// EventStore which implements no-op versions of all functions required by the interface
type StubEventStore struct{}

func (stub *StubEventStore) CountInTimeRange(name string, start, end int64) (int, error) {
	if name == "foo" {
		return 18, nil
	} else if name == "bar" {
		return 6, nil
	} else {
		return 0, nil
	}
}

func (stub *StubEventStore) Names() ([]string, error) {
	return []string{"foo", "bar", "test"}, nil
}

func (stub *StubEventStore) Put(event domain.Event) error {
	return nil
}

// EventStore which simulates an error from Put()
type StubEventStoreWithPutError struct {
	StubEventStore
}

func (stub *StubEventStoreWithPutError) Put(event domain.Event) error {
	return errors.New("error from EventStore->Put")
}

// EventStore which simulates an error from CountInTimeRange()
type StubEventStoreWithCountError struct {
	StubEventStore
}

func (stub *StubEventStoreWithCountError) CountInTimeRange(name string, start, end int64) (int, error) {
	return 0, errors.New("error from EventStore->CountInTimeRange")
}

// EventStore which simulates an error from Names()
type StubEventStoreWithNamesError struct {
	StubEventStore
}

func (stub *StubEventStoreWithNamesError) Names() ([]string, error) {
	return []string{}, errors.New("error from EventStore->Names")
}
