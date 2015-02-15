package interfaces

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

type RedisEventStore struct {
	conn *redis.Conn
}

func (repo *RedisEventStore) Put(event domain.Event) error {
	fmt.Println("RedisEventStore->Put")
	return nil
}

func NewRedisEventStore(addr, port string) (RedisEventStore, error) {
	address := fmt.Sprintf("%s:%s", addr, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return RedisEventStore{}, errors.New("Unable to connect to redis")
	}
	return RedisEventStore{conn: &conn}, nil
}
