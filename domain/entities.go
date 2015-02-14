// Package domain defines the primitive entities present in
// the events service.
package domain

type EventRepo interface {
	Store(event Event) (Event, error)
}

type Event struct {
	Name      string
	Timestamp string
}
