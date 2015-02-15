package interfaces

import (
	"log"
	"testing"

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

func TestNewRedisEventStore(t *testing.T) {
	server := startRedis("12313")
	if _, err := NewRedisEventStore("127.0.0.1", "12313"); err != nil {
		t.Fail()
	}
	stopRedis(server)
}

func TestNewRedisEventStoreUnableToConnect(t *testing.T) {
	expectedError := "Unable to connect to redis"
	_, err := NewRedisEventStore("127.0.0.1", "6379")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
