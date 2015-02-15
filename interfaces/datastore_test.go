package interfaces

import (
	"fmt"
	"log"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"

	"github.com/declantraynor/go-events-service/domain"
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
	defer stopRedis(server)

	if _, err := NewRedisEventStore("127.0.0.1", "12313"); err != nil {
		t.Fail()
	}
}

func TestNewRedisEventStoreUnableToConnect(t *testing.T) {
	expectedError := "Unable to connect to redis"
	_, err := NewRedisEventStore("127.0.0.1", "6379")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestGetEventKeyFormat(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	store, _ := NewRedisEventStore("127.0.0.1", "12313")

	expected := "event:1"
	if key, _ := store.getEventKey(); key != expected {
		t.Errorf("expected event key %q, got %q", expected, key)
	}
}

func TestGetEventKeyIncrementsByOne(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	store, _ := NewRedisEventStore("127.0.0.1", "12313")

	var key string
	for i := 1; i < 11; i++ {
		expected := fmt.Sprintf("event:%d", i)
		key, _ = store.getEventKey()
		if key != expected {
			t.Errorf("expected key %q, got %q", expected, key)
		}
	}
}

func TestPutEvent(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	store, _ := NewRedisEventStore("127.0.0.1", "12313")
	event := domain.Event{Name: "test", Timestamp: "2015-02-11T15:01:00+00:00"}

	if err := store.Put(event); err != nil {
		t.Error(err)
	}

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	name, err := redis.String(conn.Do("HGET", "event:1", "name"))
	if err != nil {
		t.Error(err)
	}
	if name != event.Name {
		t.Error(err)
	}
}

func TestPutEventError(t *testing.T) {
	server := startRedis("12313")
	store, _ := NewRedisEventStore("127.0.0.1", "12313")
	event := domain.Event{Name: "test", Timestamp: "2015-02-11T15:01:00+00:00"}

	// simulate redis connection loss
	stopRedis(server)

	if err := store.Put(event); err == nil {
		t.Errorf("expected error from RedisEventStore")
	}
}
