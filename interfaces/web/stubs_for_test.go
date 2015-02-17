package web

import (
	"errors"

	"github.com/declantraynor/go-events-service/usecases"
)

// EventInteractor which implements no-op versions of all
// functions required by the interface
type StubEventInteractor struct{}

func (interactor *StubEventInteractor) AddEvent(name, timestamp string) error {
	return nil
}

func (interactor *StubEventInteractor) CountEventsInTimeRange(from, to string) (map[string]int, error) {
	return map[string]int{
		"foo": 25,
		"bar": 43,
	}, nil
}

// EventInteractor which simulates an error from AddEvent
type StubEventInteractorWithAddError struct {
	StubEventInteractor
}

func (interactor *StubEventInteractorWithAddError) AddEvent(name, timestamp string) error {
	return errors.New("error from EventInteractor->AddEvent")
}

// EventInteractor which simulates an unspecified error from CountEventsInTimeRange
type StubEventInteractorWithCountError struct {
	StubEventInteractor
}

func (interactor *StubEventInteractorWithCountError) CountEventsInTimeRange(from, to string) (map[string]int, error) {
	return map[string]int{}, errors.New("error from EventInteractor->CountEventsInTimeRange")
}

// EventInteractor which simulates an InvalidTimestampError from CountEventsInTimeRange
type StubEventInteractorWithTimestampError struct {
	StubEventInteractor
}

func (interactor *StubEventInteractorWithTimestampError) CountEventsInTimeRange(from, to string) (map[string]int, error) {
	return map[string]int{}, usecases.InvalidTimestampError{Timestamp: from}
}

// EventInteractor which simulates an InvalidTimeRangeError from CountEventsInTimeRange
type StubEventInteractorWithTimeRangeError struct {
	StubEventInteractor
}

func (interactor *StubEventInteractorWithTimeRangeError) CountEventsInTimeRange(from, to string) (map[string]int, error) {
	return map[string]int{}, usecases.InvalidTimeRangeError{From: from, To: to}
}
