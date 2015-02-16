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

func TestNewRedisEventStore(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	if _, err := NewRedisEventStore("127.0.0.1", "12313"); err != nil {
		t.Fail()
	}
}

func TestNewRedisEventConnectionError(t *testing.T) {
	expectedError := "error connecting to redis"
	_, err := NewRedisEventStore("127.0.0.1", "6379")
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestCountInTimeRange(t *testing.T) {
	cases := []struct {
		name      string
		timestamp int64
	}{
		{"test", 1423666860},
		{"test", 1423666860},
		{"foo", 1423666860},
		{"bar", 1423666861},
		{"test", 1423666861},
		{"test", 1423666862},
		{"bar", 1423666863},
		{"foo", 1423666864},
		{"test", 1423666864},
		{"foo", 1423666870},
	}

	server := startRedis("12313")
	defer stopRedis(server)

	store, _ := NewRedisEventStore("127.0.0.1", "12313")
	for _, c := range cases {
		event := domain.Event{Name: c.name, Timestamp: c.timestamp}
		store.Put(event)
	}

	count, err := store.CountInTimeRange("test", 1423666860, 1423666870)
	if err != nil || count != 5 {
		t.Errorf("expected %d events in time range, got %d", 5, count)
	}
}

func TestCountInTimeRangeConnectionError(t *testing.T) {
	server := startRedis("12313")
	store, _ := NewRedisEventStore("127.0.0.1", "12313")

	// simulate redis connection loss
	stopRedis(server)

	if _, err := store.CountInTimeRange("test", 1423666860, 1423666870); err == nil {
		t.Fail()
	}
}

func TestNames(t *testing.T) {
	cases := []struct {
		name      string
		timestamp int64
	}{
		{"test", 1423666860},
		{"test", 1423666860},
		{"foo", 1423666860},
		{"bar", 1423666861},
		{"test", 1423666861},
		{"test", 1423666862},
		{"bar", 1423666863},
		{"foo", 1423666864},
		{"test", 1423666864},
		{"foo", 1423666870},
	}

	server := startRedis("12313")
	defer stopRedis(server)

	store, _ := NewRedisEventStore("127.0.0.1", "12313")
	for _, c := range cases {
		event := domain.Event{Name: c.name, Timestamp: c.timestamp}
		store.Put(event)
	}

	expected := []string{"test", "foo", "bar"}
	names, _ := store.Names()

	for _, name := range expected {
		if !stringInSlice(name, names) {
			t.Errorf("expected value %q not present in names")
		}
	}
}

func TestPut(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &PassingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: 1423666860}

	if err := store.Put(event); err != nil {
		t.Fail()
	}

	// event name is stored
	nameStored, err := redis.Int(conn.Do("SISMEMBER", "event_names", event.Name))
	if err != nil || nameStored != 1 {
		t.Errorf("event name not stored")
	}

	// event entity is stored
	name, err := redis.String(conn.Do("HGET", "event:1", "name"))
	if err != nil || name != event.Name {
		t.Error("event entity not stored")
	}

	// sorted event index (events by name and timestamp) is created
	numSortedEvents, err := redis.Int(conn.Do("ZCARD", "events:test:by-timestamp"))
	if err != nil || numSortedEvents != 1 {
		t.Error("sorted event index not created")
	}

}

func TestNamesConnectionError(t *testing.T) {
	server := startRedis("12313")
	store, _ := NewRedisEventStore("127.0.0.1", "12313")

	// simulate redis connection loss
	stopRedis(server)

	if _, err := store.Names(); err == nil {
		t.Fail()
	}
}

func TestPutConnectionError(t *testing.T) {
	server := startRedis("12313")

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &PassingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: 1423666860}

	// simulate redis connection loss
	stopRedis(server)

	if err := store.Put(event); err == nil {
		t.Fail()
	}
}

func TestPutIdGeneratorError(t *testing.T) {
	server := startRedis("12313")
	defer stopRedis(server)

	conn, _ := redis.Dial("tcp", "127.0.0.1:12313")
	defer conn.Close()

	store := RedisEventStore{conn: conn, idgen: &FailingIdGenerator{}}
	event := domain.Event{Name: "test", Timestamp: 1423666860}

	if err := store.Put(event); err == nil {
		t.Fail()
	}
}

func TestSanitizeName(t *testing.T) {
	cases := []struct {
		input, expect string
	}{
		{"name", "name"},
		{"name-with-dashes", "name-with-dashes"},
		{" name-with-leading-space", "name-with-leading-space"},
		{" name-with-leading-and-trailing-space ", "name-with-leading-and-trailing-space"},
		{"name with internal spaces", "name-with-internal-spaces"},
		{"  name with multiple spaces  ", "name-with-multiple-spaces"},
	}

	for _, c := range cases {
		if result := sanitizeName(c.input); result != c.expect {
			t.Errorf("expected %q, got %q", c.expect, result)
		}
	}
}
