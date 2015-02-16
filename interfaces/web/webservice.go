package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type EventInteractor interface {
	Add(name, timestamp string) error
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

func (handler *WebService) Create(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		handler.RenderJSON(res, ErrorResource{Error: "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	event := EventResource{}
	err := json.Unmarshal(body, &event)
	if err != nil {
		handler.RenderJSON(res, ErrorResource{Error: "request JSON is invalid"}, http.StatusBadRequest)
		return
	}

	if err := handler.EventInteractor.Add(event.Name, event.Timestamp); err != nil {
		handler.RenderJSON(res, ErrorResource{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	handler.RenderJSON(res, "", http.StatusCreated)
}

func (handler *WebService) RenderJSON(res http.ResponseWriter, resource interface{}, status int) {
	responseBody, _ := json.Marshal(resource)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(status)
	res.Write(responseBody)
}
