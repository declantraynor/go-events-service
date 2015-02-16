package usecases

import (
	"testing"
)

func TestAddEventSucceeds(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	if err := interactor.AddEvent("test-event", "2015-02-11T15:01:00+00:00"); err != nil {
		t.Fail()
	}
}

func TestAddEventNonISOTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	err := interactor.AddEvent("test-event", "2015/02/01 15:01")

	if err == nil || err.Error() != "timestamps must conform to ISO8601" {
		t.Fail()
	}
}

func TestAddEventNonUTCTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(PassingEventStore)}
	err := interactor.AddEvent("test-event", "2015-02-11T15:01:00-05:00")

	if err == nil || err.Error() != "timestamps must be UTC" {
		t.Fail()
	}
}

func TestAddEventStorageError(t *testing.T) {
	interactor := EventInteractor{Store: new(FailingEventStore)}
	if err := interactor.AddEvent("test-event", "2015-02-11T15:01:00+00:00"); err == nil {
		t.Fail()
	}
}
