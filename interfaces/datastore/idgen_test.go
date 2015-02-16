package datastore

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestNextIncrementsByOne(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	gen := RedisIdGenerator{conn: conn, name: "test"}

	var expectedId int64
	for expectedId = 1; expectedId < 11; expectedId++ {
		if id, _ := gen.Next(); id != expectedId {
			t.Errorf("expected ID %d, got %d", expectedId, id)
		}
	}
}

func TestNextEncountersConnectionError(t *testing.T) {
	server := startRedis("12313")

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	gen := RedisIdGenerator{conn: conn, name: "test"}

	// simulate redis connection loss
	stopRedis(server)

	id, err := gen.Next()
	if id != 0 || err == nil {
		t.Fail()
	}
}
