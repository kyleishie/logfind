package logfind

import (
	"fmt"
	"github.com/kyleishie/logfind/pkg/logfind/reader"
	"io"
	"strings"
	"time"
)

// Finder represents an object that is capable of using finderOptions to query an input log stream.
//
// Note: It is up to the concrete implementation of Finder as to how it ingests a log stream, e.g.,
// as a file, as a slice of strings, as a reader, etc.
type Finder interface {

	// Find applies the given opts to each event in a log stream to find match records.
	Find(opts ...FinderOptionFunc) (count int, events []string, err error)
}

type defaultFinder struct {
	r reader.Reader
}

func NewFinder(r reader.Reader) Finder {
	return &defaultFinder{
		r: r,
	}
}

func (f *defaultFinder) Find(opts ...FinderOptionFunc) (count int, events []string, err error) {
	// Clean up on error
	defer func() {
		if err != nil && err != io.EOF {
			count = 0
			events = nil
		}
	}()

	options, err := newFindOptions(opts...)
	if err != nil {
		return
	}

	var countMap map[string]bool
	if options.cc != Event {
		countMap = make(map[string]bool)
		defer func() {
			// only use the countmap length when countConcern is something other than Event
			count = len(countMap)
		}()
	}

	for {
		event, err := f.r.Read()
		if err != nil {
			break
		}

		timestampMatch := true
		usernameMatch := true
		operationMatch := true
		minSizeMatch := true
		maxSizeMatch := true

		timestamp, err := event.Timestamp()
		if err != nil {
			break
		}
		if options.minTime != nil && options.maxTime != nil {
			timestampMatch = timestamp.After(*options.minTime) && timestamp.Before(*options.maxTime)
		}

		username, err := event.Username()
		if err != nil {
			break
		}
		if options.username != nil {
			usernameMatch = strings.EqualFold(username, *options.username)
		}

		operation, err := event.Operation()
		if err != nil {
			break
		}
		if options.operation != nil {
			operationMatch = strings.EqualFold(operation, *options.operation)
		}

		size, err := event.Size()
		if err != nil {
			break
		}
		if options.minSize != nil {
			minSizeMatch = size >= *options.minSize
		}

		if options.maxSize != nil {
			maxSizeMatch = size <= *options.maxSize
		}

		if timestampMatch && usernameMatch && operationMatch && minSizeMatch && maxSizeMatch {
			/// Match Found
			events = append(events, fmt.Sprintf("%s %s %s %d", timestamp.Format(time.UnixDate), username, operation, size))

			/// Count the uniqueness of the match based on countConcern
			switch options.cc {
			case Event:
				count++
			case Operation:
				countMap[operation] = true
			case User:
				countMap[username] = true
			}
		}

	}

	return
}
