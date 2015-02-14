package interfaces

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/declantraynor/go-events-service/domain"
	"github.com/declantraynor/go-events-service/usecases"
)

type StubEventRepo struct{}

func (stub *StubEventRepo) Store(event domain.Event) (domain.Event, error) {
	return event, nil
}

func TestCreateEventSucceeds(t *testing.T) {
	interactor := usecases.EventInteractor{Repo: new(StubEventRepo)}
	handler := new(WebServiceHandler)
	handler.eventInteractor = interactor

	requestBody := strings.NewReader(`{"name": "test", "timestamp": "2015-02-11T15:01:00+00:00"}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	handler.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.Code)
	}
}

func TestCreateEventRejectsInvalidJSON(t *testing.T) {
	interactor := usecases.EventInteractor{Repo: new(StubEventRepo)}
	handler := new(WebServiceHandler)
	handler.eventInteractor = interactor

	requestBody := strings.NewReader(`{"invalid": json}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.Code)
	}
}
