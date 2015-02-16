package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type PassingEventInteractor struct{}
type FailingEventInteractor struct{}

func (interactor *PassingEventInteractor) Add(name, timestamp string) error {
	return nil
}

func (interactor *FailingEventInteractor) Add(name, timestamp string) error {
	return errors.New("error from event interactor Add")
}

func TestCreateSucceeds(t *testing.T) {
	service := WebService{EventInteractor: new(PassingEventInteractor)}

	requestBody := strings.NewReader(`{"name": "test", "timestamp": "2015-02-11T15:01:00+00:00"}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	service.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.Code)
	}
}

func TestCreateRejectsInvalidHTTPMethods(t *testing.T) {
	methods := []string{"GET", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	service := WebService{EventInteractor: new(PassingEventInteractor)}

	for _, m := range methods {
		request, _ := http.NewRequest(m, "http://example.com/events", nil)
		response := httptest.NewRecorder()
		service.Create(response, request)

		expectedCode := http.StatusMethodNotAllowed
		if response.Code != expectedCode {
			t.Errorf("expected response code %d, got %d", expectedCode, response.Code)
		}
	}
}

func TestCreateRejectsInvalidJSON(t *testing.T) {
	service := WebService{EventInteractor: new(PassingEventInteractor)}

	requestBody := strings.NewReader(`{"invalid": json}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	service.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.Code)
	}

	expectedResponseBody := `{"error":"request JSON is invalid"}`
	if response.Body.String() != expectedResponseBody {
		t.Errorf("expected response body %q, got %q", expectedResponseBody, response.Body.String())
	}
}

func TestCreateHandlesInteractorError(t *testing.T) {
	service := WebService{EventInteractor: new(FailingEventInteractor)}

	requestBody := strings.NewReader(`{"name": "test", "timestamp": "2015-02-11T15:01:00-05:00"}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	service.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", http.StatusCreated, response.Code)
	}

	expectedResponseBody := `{"error":"error from event interactor Add"}`
	if response.Body.String() != expectedResponseBody {
		t.Errorf("expected response body %q, got %q", expectedResponseBody, response.Body.String())
	}
}
