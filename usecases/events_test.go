package usecases

import (
	"reflect"
	"testing"
)

func TestAddEvent(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}

	if err := interactor.AddEvent("test-event", "2015-02-11T15:01:00+00:00"); err != nil {
		t.Error("EventInteractor.AddEvent returned an unexpected error")
	}
}

func TestAddEventNonISOTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	err := interactor.AddEvent("test-event", "2015/02/01 15:01")

	if err, ok := err.(InvalidTimestampError); !ok || err.NotISO8601 == false {
		t.Errorf("expected InvalidTimestampError, got %T", err)
	}

	expectedErrorFormat := `2015/02/01 15:01 does not conform to ISO8601`
	if err.Error() != expectedErrorFormat {
		t.Error("InvalidTimestampError format is wrong")
	}
}

func TestAddEventNonUTCTimestamp(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	err := interactor.AddEvent("test-event", "2015-02-11T15:01:00-05:00")

	if err, ok := err.(InvalidTimestampError); !ok || err.NotUTC == false {
		t.Errorf("expected InvalidTimestampError, got %T", err)
	}

	expectedErrorFormat := `2015-02-11T15:01:00-05:00 is not UTC`
	if err.Error() != expectedErrorFormat {
		t.Error("InvalidTimestampError format is wrong")
	}
}

func TestAddEventStorageError(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStoreWithPutError)}

	if err := interactor.AddEvent("test-event", "2015-02-11T15:01:00+00:00"); err == nil {
		t.Error("expected error from Store.Put")
	}
}

func TestCountEventsInTimeRange(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	counts, err := interactor.CountEventsInTimeRange("2015-01-01T13:23:00+00:00", "2015-01-01T13:23:59+00:00")

	if err != nil {
		t.Error("EventInteractor.CountEventsInTimeRange returned unexpected error")
	}

	// expect all non-zero count values returned by StubEventStore
	expected := map[string]int{
		"foo": 18,
		"bar": 6,
	}
	if !reflect.DeepEqual(counts, expected) {
		t.Errorf("expected %v, got %v", expected, counts)
	}
}

func TestCountEventsInTimeRangeInvalidFrom(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	_, err := interactor.CountEventsInTimeRange("2015/01/01 13:23:00", "2015-01-01T13:23:59+00:00")

	if _, ok := err.(InvalidTimestampError); !ok {
		t.Errorf("expected InvalidTimestampError, got %T", err)
	}
}

func TestCountEventsInTimeRangeInvalidTo(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	_, err := interactor.CountEventsInTimeRange("2015-01-01T13:23:00+00:00", "2015/01/01 13:23:59")

	if _, ok := err.(InvalidTimestampError); !ok {
		t.Errorf("expected InvalidTimestampError, got %T", err)
	}
}

func TestCountEventsInTimeRangeInvalidRange(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStore)}
	_, err := interactor.CountEventsInTimeRange("2015-01-01T13:29:00+00:00", "2015-01-01T13:20:00+00:00")

	if _, ok := err.(InvalidTimeRangeError); !ok {
		t.Errorf("expected InvalidTimeRangeError, got %T", err)
	}

	expectedErrorFormat := `2015-01-01T13:29:00+00:00 is later than 2015-01-01T13:20:00+00:00`
	if err.Error() != expectedErrorFormat {
		t.Error("InvalidTimeRangeError format is wrong")
	}
}

func TestCountEventsInTimeRangeEventStoreNamesError(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStoreWithNamesError)}
	_, err := interactor.CountEventsInTimeRange("2015-01-01T13:23:00+00:00", "2015-01-01T13:23:59+00:00")

	if err == nil {
		t.Error("expected error from Store.Names")
	}
}

func TestCountEventsInTimeRangeEventStoreCountError(t *testing.T) {
	interactor := EventInteractor{Store: new(StubEventStoreWithCountError)}
	_, err := interactor.CountEventsInTimeRange("2015-01-01T13:23:00+00:00", "2015-01-01T13:23:59+00:00")

	if err == nil {
		t.Error("expected error from Store.CountInTimeRange")
	}
}
