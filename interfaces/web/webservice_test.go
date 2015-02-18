package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractor)}

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
	service := WebService{EventInteractor: new(StubEventInteractor)}

	for _, m := range methods {
		request, _ := http.NewRequest(m, "http://example.com/events", nil)
		response := httptest.NewRecorder()
		service.Create(response, request)

		expectedResponseCode := http.StatusMethodNotAllowed
		if response.Code != expectedResponseCode {
			t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
		}
	}
}

func TestCreateRejectsInvalidJSON(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractor)}

	requestBody := strings.NewReader(`{"invalid": json}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	service.Create(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCreateEventInteractorError(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractorWithAddError)}

	requestBody := strings.NewReader(`{"name": "test", "timestamp": "2015-02-11T15:01:00-05:00"}`)
	request, _ := http.NewRequest("POST", "http://example.com/events", requestBody)

	response := httptest.NewRecorder()
	service.Create(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != http.StatusBadRequest {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCount(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractor)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?from=2015-02-11T15:01:00+00:00&to=2015-02-11T15:01:59+00:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusOK
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}

	// expect values returned by StubEventInteractor
	expectedResponse := map[string]int{
		"foo": 25,
		"bar": 43,
	}

	receivedResponse := make(map[string]int)
	json.Unmarshal(response.Body.Bytes(), &receivedResponse)

	if !reflect.DeepEqual(receivedResponse, expectedResponse) {
		t.Errorf("response is incorrect, got: %s", response.Body.String())
	}
}

func TestCountRejectsInvalidRequestMethods(t *testing.T) {
	methods := []string{"POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	service := WebService{EventInteractor: new(StubEventInteractor)}

	for _, m := range methods {
		request, _ := http.NewRequest(m, "http://example.com/events/count", nil)
		response := httptest.NewRecorder()
		service.Count(response, request)

		expectedResponseCode := http.StatusMethodNotAllowed
		if response.Code != expectedResponseCode {
			t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
		}
	}
}

func TestCountMissingFromParameter(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractor)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?to=2015-02-11T15:01:59+00:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCountMissingToParameter(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractor)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?from=2015-02-11T15:01:00+00:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCountInvalidTimestampError(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractorWithTimestampError)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?from=2015-02-11T15:01:00-05:00&to=2015-02-11T15:01:59-05:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCountInvalidTimeRangeError(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractorWithTimeRangeError)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?from=2015-02-15T00:01:00+00:00&to=2015-02-11T15:01:59+00:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusBadRequest
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}

func TestCountGenericInteractorError(t *testing.T) {
	service := WebService{EventInteractor: new(StubEventInteractorWithCountError)}
	request, _ := http.NewRequest(
		"GET",
		"http://example.com/events/count?from=2015-02-15T00:01:00+00:00&to=2015-02-15T15:01:59+00:00",
		nil)

	response := httptest.NewRecorder()
	service.Count(response, request)

	expectedResponseCode := http.StatusInternalServerError
	if response.Code != expectedResponseCode {
		t.Errorf("expected response code %d, got %d", expectedResponseCode, response.Code)
	}
}
