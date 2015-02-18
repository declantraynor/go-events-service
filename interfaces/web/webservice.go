// Package web implements a number of handlers which expose the the service's functionality via HTTP.
package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/declantraynor/go-events-service/usecases"
)

type EventInteractor interface {
	AddEvent(name, timestamp string) error
	CountEventsInTimeRange(from, to string) (map[string]int, error)
}

type EventResource struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

type ErrorResource struct {
	Error string `json:"error"`
}

type WebService struct {
	EventInteractor EventInteractor
}

func (service *WebService) Create(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		service.RenderJSON(
			res,
			ErrorResource{Error: "Method Not Allowed"},
			http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	event := EventResource{}
	err := json.Unmarshal(body, &event)
	if err != nil {
		service.RenderJSON(
			res,
			ErrorResource{Error: "Request JSON is invalid"},
			http.StatusBadRequest)
		return
	}

	if err := service.EventInteractor.AddEvent(event.Name, event.Timestamp); err != nil {
		service.RenderJSON(
			res,
			ErrorResource{Error: err.Error()},
			http.StatusBadRequest)
		return
	}

	service.RenderJSON(res, map[string]string{}, http.StatusCreated)
}

func (service *WebService) Count(res http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		service.RenderJSON(
			res,
			ErrorResource{Error: "Method Not Allowed"},
			http.StatusMethodNotAllowed)
		return
	}

	// FormValue will parse out any `+` symbols in query params,
	// so we need to put them back in to get the true timestamp
	// values passed in the URL
	from := strings.Replace(req.FormValue("from"), " ", "+", -1)
	to := strings.Replace(req.FormValue("to"), " ", "+", -1)

	if from == "" {
		service.RenderJSON(
			res,
			ErrorResource{Error: `Missing required parameter "from"`},
			http.StatusBadRequest)
		return
	}

	if to == "" {
		service.RenderJSON(
			res,
			ErrorResource{Error: `Missing required parameter "to"`},
			http.StatusBadRequest)
		return
	}

	counts, err := service.EventInteractor.CountEventsInTimeRange(from, to)
	if err != nil {
		if e, ok := err.(usecases.InvalidTimestampError); ok {
			service.RenderJSON(
				res,
				ErrorResource{Error: e.Error()},
				http.StatusBadRequest)
			return
		}
		if e, ok := err.(usecases.InvalidTimeRangeError); ok {
			service.RenderJSON(
				res,
				ErrorResource{Error: e.Error()},
				http.StatusBadRequest)
			return
		}
		service.RenderJSON(
			res,
			ErrorResource{Error: "Internal Server Error"},
			http.StatusInternalServerError)
		return
	}

	service.RenderJSON(res, counts, http.StatusOK)
}

func (service *WebService) RenderJSON(res http.ResponseWriter, resource interface{}, status int) {
	responseBody, _ := json.Marshal(resource)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(status)
	res.Write(responseBody)
}
