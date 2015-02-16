package datastore

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

type RedisEventStore struct {
	conn  redis.Conn
	idgen IdGenerator
}

func (store *RedisEventStore) store(key string, event domain.Event) error {
	_, err := store.conn.Do("HMSET", key, "name", event.Name, "timestamp", event.Timestamp)
	if err != nil {
		return errors.New("error storing event")
	}
	return nil
}

func (store *RedisEventStore) Put(event domain.Event) error {
	id, err := store.idgen.Next()
	if err != nil {
		return errors.New("error generating event ID")
	}

	key := fmt.Sprintf("event:%d", id)
	return store.store(key, event)
}

func NewRedisEventStore(addr, port string) (RedisEventStore, error) {
	address := fmt.Sprintf("%s:%s", addr, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return RedisEventStore{}, errors.New("error connecting to redis")
	}
	idgen := RedisIdGenerator{conn: conn, name: "next_event_id"}
	return RedisEventStore{conn: conn, idgen: &idgen}, nil
}
