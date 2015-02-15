package main

import (
	"log"
	"net/http"
	"os"

	"github.com/declantraynor/go-events-service/interfaces"
	"github.com/declantraynor/go-events-service/usecases"
)

func main() {

	redisAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	eventStore, err := interfaces.NewRedisEventStore(redisAddr, redisPort)
	if err != nil {
		log.Fatal(err.Error())
	}

	eventInteractor := usecases.EventInteractor{Store: &eventStore}
	webservice := interfaces.WebService{EventInteractor: &eventInteractor}

	http.HandleFunc("/events", func(res http.ResponseWriter, req *http.Request) {
		webservice.Create(res, req)
	})
	http.ListenAndServe(":5000", nil)
}
