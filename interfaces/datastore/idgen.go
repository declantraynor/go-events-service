package datastore

import (
	"github.com/garyburd/redigo/redis"
)

type IdGenerator interface {
	Next() (int64, error)
}

type RedisIdGenerator struct {
	conn redis.Conn
	name string
}

func (gen *RedisIdGenerator) Next() (int64, error) {
	id, err := redis.Int64(gen.conn.Do("INCR", gen.name))
	if err != nil {
		return 0, err
	}
	return id, nil
}
