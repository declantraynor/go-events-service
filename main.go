package main

import (
	"net/http"

	"github.com/declantraynor/go-events-service/domain"
	"github.com/declantraynor/go-events-service/interfaces"
	"github.com/declantraynor/go-events-service/usecases"
)

type StubEventRepo struct{}

func (stub *StubEventRepo) Store(event domain.Event) (domain.Event, error) {
	return domain.Event{}, nil
}

func main() {
	eventInteractor := new(usecases.EventInteractor)
	eventInteractor.Repo = new(StubEventRepo)
	webservice := interfaces.WebService{}
	webservice.EventInteractor = eventInteractor

	http.HandleFunc("/events", func(res http.ResponseWriter, req *http.Request) {
		webservice.Create(res, req)
	})
	http.ListenAndServe(":5000", nil)
}
