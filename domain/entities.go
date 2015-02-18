// Package domain defines the primitive entities present in the events service.
package domain

type EventStore interface {
	CountInTimeRange(name string, start, end int64) (int, error)
	Names() ([]string, error)
	Put(event Event) error
}

type Event struct {
	Name      string
	Timestamp int64
}
