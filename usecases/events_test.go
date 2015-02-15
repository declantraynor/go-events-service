package usecases

import (
	"testing"

	"github.com/declantraynor/go-events-service/domain"
)

type StubEventStore struct{}

func (stub *StubEventStore) Put(event domain.Event) error {
	return nil
}

func TestAddEventSucceeds(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	if err := interactor.Add("test-event", "2015-02-11T15:01:00+00:00"); err != nil {
		t.Fail()
	}
}

func TestAddEventChecksForISOTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	err := interactor.Add("test-event", "2015/02/01 15:01")

	if err == nil || err.Error() != "timestamps must conform to ISO8601" {
		t.Fail()
	}
}

func TestAddEventChecksForUTCTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	err := interactor.Add("test-event", "2015-02-11T15:01:00-05:00")

	if err == nil || err.Error() != "timestamps must be UTC" {
		t.Fail()
	}
}
