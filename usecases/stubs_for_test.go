package usecases

import (
	"errors"

	"github.com/declantraynor/go-events-service/domain"
)

// PassingEventStore
type PassingEventStore struct{}

func (stub *PassingEventStore) CountInTimeRange(name string, start, end int64) (int, error) {
	return 0, nil
}

func (stub *PassingEventStore) Names() ([]string, error) {
	return []string{}, nil
}

func (stub *PassingEventStore) Put(event domain.Event) error {
	return nil
}

// FailingEventStore
type FailingEventStore struct{}

func (stub *FailingEventStore) CountInTimeRange(name string, start, end int64) (int, error) {
	return 0, errors.New("error from EventStore->CountInTimeRange")
}

func (stub *FailingEventStore) Names() ([]string, error) {
	return []string{}, errors.New("error from EventStore->Names")
}

func (stub *FailingEventStore) Put(event domain.Event) error {
	return errors.New("error from EventStore->Put")
}
