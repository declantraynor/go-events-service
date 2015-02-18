package main

import (
	"log"
	"net/http"
	"os"

	"github.com/declantraynor/go-events-service/interfaces/datastore"
	"github.com/declantraynor/go-events-service/interfaces/web"
	"github.com/declantraynor/go-events-service/usecases"
)

func run(serve func(webservice *web.WebService)) error {
	redisAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	eventStore, err := datastore.NewRedisEventStore(redisAddr, redisPort)
	if err != nil {
		return err
	}

	eventInteractor := usecases.EventInteractor{Store: &eventStore}
	webservice := web.WebService{EventInteractor: &eventInteractor}

	serve(&webservice)
	return nil
}

func serve(webservice *web.WebService) {
	http.HandleFunc("/events", webservice.Create)
	http.HandleFunc("/events/count", webservice.Count)
	http.ListenAndServe(":5000", nil)
}

func main() {
	if err := run(serve); err != nil {
		log.Fatal(err.Error())
	}
}
