package datastore

import (
	"log"

	"github.com/stvp/tempredis"
)

func startRedis(port string) *tempredis.Server {
	server, err := tempredis.Start(
		tempredis.Config{
			"port": port,
		},
	)
	if err != nil {
		log.Fatal("Unable to start tempredis for test")
	}
	return server
}

func stopRedis(server *tempredis.Server) {
	err := server.Kill()
	if err != nil {
		log.Fatal("Problem killing tempredis server during test")
	}
}
