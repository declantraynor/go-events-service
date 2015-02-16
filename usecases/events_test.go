package usecases

import (
	"testing"
)

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
