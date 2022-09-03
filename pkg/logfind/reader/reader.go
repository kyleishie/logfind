package reader

import (
	"time"
)

// Reader represents an object that can parse events from a log stream.
type Reader interface {
	Read() (event Event, err error)
}

// Event represents an event from a log stream.
//
// Note: This is intentionally not a concrete type in order to prevent the need for
// type conversions within Reader implementations as the result of, for instance,
// using struct tags to unmarshal log events.
type Event interface {
	Timestamp() (time.Time, error)
	Username() (string, error)
	Operation() (string, error)
	Size() (int, error)
}
