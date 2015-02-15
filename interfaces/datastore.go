package interfaces

import (
	"fmt"

	"github.com/garyburd/redigo/redis"

	"github.com/declantraynor/go-events-service/domain"
)

type RedisEventRepo struct {
	Conn *redis.Conn
}

func (repo *RedisEventRepo) Store(event domain.Event) (domain.Event, error) {
	fmt.Print("RedisEventRepo->Store")
	return event, nil
}
