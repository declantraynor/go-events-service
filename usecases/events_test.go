package usecases

import (
	"errors"
	"testing"

	"github.com/declantraynor/go-events-service/domain"
)

type PassingEventStore struct{}

func (stub *PassingEventStore) Put(event domain.Event) error {
	return nil
}

type FailingEventStore struct{}

func (stub *FailingEventStore) Put(event domain.Event) error {
	return errors.New("error from EventStore->Put")
}

func TestAddSucceeds(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	if err := interactor.Add("test-event", "2015-02-11T15:01:00+00:00"); err != nil {
		t.Fail()
	}
}

func TestAddChecksForISOTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	err := interactor.Add("test-event", "2015/02/01 15:01")

	if err == nil || err.Error() != "timestamps must conform to ISO8601" {
		t.Fail()
	}
}

func TestAddChecksForUTCTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	err := interactor.Add("test-event", "2015-02-11T15:01:00-05:00")

	if err == nil || err.Error() != "timestamps must be UTC" {
		t.Fail()
	}
}

func TestAddEncountersEventStoreError(t *testing.T) {
	interactor := EventInteractor{Store: new(FailingEventStore)}
	if err := interactor.Add("test-event", "2015-02-11T15:01:00+00:00"); err == nil {
		t.Fail()
	}
}
