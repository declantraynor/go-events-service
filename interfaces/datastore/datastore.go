// Package datastore implements an interface to a redis datastore
package datastore

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

func sanitizeName(name string) string {
	re := regexp.MustCompile("\\s+")
	return re.ReplaceAllString(strings.TrimSpace(name), "-")
}

type RedisEventStore struct {
	conn  redis.Conn
	idgen IdGenerator
}

// CountInTimeRange returns an integer count of all events with a given name
// and timestamp between `start` and `end`, as well as any error encountered.
func (store *RedisEventStore) CountInTimeRange(name string, start, end int64) (int, error) {
	index := fmt.Sprintf("events:%s:by-timestamp", sanitizeName(name))
	count, err := redis.Int(store.conn.Do("ZCOUNT", index, start, end))
	if err != nil {
		return 0, errors.New("error getting event count")
	}
	return count, nil
}

// Names returns a string slice containing all previously stored event names,
// as well as any error encountered.
func (store *RedisEventStore) Names() ([]string, error) {
	names, err := redis.Strings(store.conn.Do("SMEMBERS", "event_names"))
	if err != nil {
		return []string{}, errors.New("error getting event names")
	}
	return names, nil
}

// Put stores a new event in redis, returning any error encountered.
func (store *RedisEventStore) Put(event domain.Event) error {
	id, err := store.idgen.Next()
	if err != nil {
		return errors.New("error generating event ID")
	}

	key := fmt.Sprintf("event:%d", id)
	return store.store(key, event)
}

func (store *RedisEventStore) store(key string, event domain.Event) error {
	index := fmt.Sprintf("events:%s:by-timestamp", sanitizeName(event.Name))

	// storing an event triggers a redis transaction comprising multiple operations
	store.conn.Send("MULTI")

	// add the event name to a set of all known event names (will do nothing if name already exists)
	store.conn.Send("SADD", "event_names", event.Name)

	// store the event data in a hash, uniquely identified by `key`
	store.conn.Send("HMSET", key, "name", event.Name, "timestamp", event.Timestamp)

	// add the event key to a sorted set of events with the same name,
	// sorted by timestamp to allow for efficient range queries
	store.conn.Send("ZADD", index, event.Timestamp, key)

	if _, err := store.conn.Do("EXEC"); err != nil {
		return errors.New("error storing event")
	}
	return nil
}

// NewRedisEventStore opens a TCP connection to a redis server at the given
// address and port. It returns an intialised RedisEventStore struct as well
// as any error encountered.
func NewRedisEventStore(addr, port string) (RedisEventStore, error) {
	address := fmt.Sprintf("%s:%s", addr, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return RedisEventStore{}, errors.New("error connecting to redis")
	}
	idgen := RedisIdGenerator{conn: conn, name: "next_event_id"}
	return RedisEventStore{conn: conn, idgen: &idgen}, nil
}
