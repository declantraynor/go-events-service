package datastore

import (
	"errors"
	"testing"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

type PassingIdGenerator struct{}

func (stub *PassingIdGenerator) Next() (int64, error) {
	return 1, nil
}

type FailingIdGenerator struct{}

func (stub *FailingIdGenerator) Next() (int64, error) {
	return 0, errors.New("error from IdGenerator->Next")
}

func TestPutSucceeds(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &PassingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: "2015-02-11T15:01:00+00:00"}

	if err := store.Put(event); err != nil {
		t.Error(err)
	}

	name, err := redis.String(conn.Do("HGET", "event:1", "name"))
	if err != nil || name != event.Name {
		t.Error(err)
	}
}

func TestPutEncountersConnectionError(t *testing.T) {
	server := startRedis("12313")

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &PassingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: "2015-02-11T15:01:00+00:00"}

	// simulate redis connection loss
	stopRedis(server)

	if err := store.Put(event); err == nil {
		t.Fail()
	}
}

func TestPutEncountersIdGeneratorError(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &FailingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: "2015-02-11T15:01:00+00:00"}

	if err := store.Put(event); err == nil {
		t.Fail()
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
	expectedError := "error connecting to redis"
	_, err := NewRedisEventStore("127.0.0.1", "6379")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
