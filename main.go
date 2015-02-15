package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/interfaces"
	"github.com/declantraynor/go-events-service/usecases"
)

func main() {

	redisAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redisPort := os.Getenv("REDIS_PORT_6379_TCP_PORT")

	redisConn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisAddr, redisPort))
	if err != nil {
		log.Fatal("Unable to connect to redis")
	}

	eventRepo := interfaces.RedisEventRepo{Conn: &redisConn}
	eventInteractor := usecases.EventInteractor{Repo: &eventRepo}
	webservice := interfaces.WebService{EventInteractor: &eventInteractor}

	http.HandleFunc("/events", func(res http.ResponseWriter, req *http.Request) {
		webservice.Create(res, req)
	})
	http.ListenAndServe(":5000", nil)
}
