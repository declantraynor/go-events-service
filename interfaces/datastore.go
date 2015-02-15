package interfaces

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

type RedisEventStore struct {
	conn redis.Conn
}

func (store *RedisEventStore) getEventKey() (string, error) {
	id, err := redis.Int(store.conn.Do("INCR", "next_event_id"))
	if err != nil {
		return "", errors.New("Unable to generate event ID")
	}
	return fmt.Sprintf("event:%d", id), nil
}

func (store *RedisEventStore) Put(event domain.Event) error {

	key, keyError := store.getEventKey()
	if keyError != nil {
		return keyError
	}

	_, err := store.conn.Do("HMSET", key, "name", event.Name, "timestamp", event.Timestamp)
	if err != nil {
		return errors.New("Unable to store event")
	}

	return nil
}

func NewRedisEventStore(addr, port string) (RedisEventStore, error) {
	address := fmt.Sprintf("%s:%s", addr, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return RedisEventStore{}, errors.New("Unable to connect to redis")
	}
	return RedisEventStore{conn: conn}, nil
}
