package usecases

import (
	"testing"

	"github.com/declantraynor/go-events-service/domain"
)

type StubEventRepo struct{}

func (stub *StubEventRepo) Store(event domain.Event) (domain.Event, error) {
	return event, nil
}

func TestAddEventSucceeds(t *testing.T) {
	interactor := EventInteractor{Repo: new(StubEventRepo)}
	if _, err := interactor.Add("test-event", "2015-02-11T15:01:00+00:00"); err != nil {
		t.Fail()
	}
}

func TestAddEventChecksForISOTimestamp(t *testing.T) {
	interactor := EventInteractor{Repo: new(StubEventRepo)}
	_, err := interactor.Add("test-event", "2015/02/01 15:01")

	if err == nil || err.Error() != "timestamps must conform to ISO8601" {
		t.Fail()
	}
}

func TestAddEventChecksForUTCTimestamp(t *testing.T) {
	interactor := EventInteractor{Repo: new(StubEventRepo)}
	_, err := interactor.Add("test-event", "2015-02-11T15:01:00-05:00")

	if err == nil || err.Error() != "timestamps must be UTC" {
		t.Fail()
	}
}

func TestAddEventSanitizesName(t *testing.T) {
	interactor := EventInteractor{Repo: new(StubEventRepo)}
	event, err := interactor.Add("  test spaces in name   ", "2015-02-11T15:01:00+00:00")

	if err != nil {
		t.Fail()
	}

	if event.Name != "test-spaces-in-name" {
		t.Fail()
	}
}
