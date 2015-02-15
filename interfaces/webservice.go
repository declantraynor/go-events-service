package interfaces

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/declantraynor/go-events-service/usecases"
)

type EventResource struct {
	Name      string `json:"name"`
	Timestamp string `json:timestamp`
}

type WebServiceHandler struct {
	eventInteractor usecases.EventInteractor
}

func (handler *WebServiceHandler) Create(response http.ResponseWriter, request *http.Request) {

	defer request.Body.Close()
	body, _ := ioutil.ReadAll(request.Body)

	event := EventResource{}
	err := json.Unmarshal(body, &event)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}

	if err := handler.eventInteractor.Add(event.Name, event.Timestamp); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}

	response.WriteHeader(http.StatusCreated)
}
