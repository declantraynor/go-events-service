// Package domain defines the primitive entities present in
// the events service.
package domain

type EventStore interface {
	Put(event Event) error
}

type Event struct {
	Name      string
	Timestamp int64
}
