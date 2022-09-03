package logfind

import (
	"fmt"
	"github.com/kyleishie/logfind/pkg/logfind/reader"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

type mockEvent struct {
	timestamp time.Time
	username  string
	operation string
	size      int
}

func (m mockEvent) Timestamp() (time.Time, error) {
	return m.timestamp, nil
}

func (m mockEvent) Username() (string, error) {
	return m.username, nil
}

func (m mockEvent) Operation() (string, error) {
	return m.operation, nil
}

func (m mockEvent) Size() (int, error) {
	return m.size, nil
}

type mockReader struct {
	events []mockEvent
	i      int
}

func (m *mockReader) Read() (e reader.Event, err error) {
	if m.i >= len(m.events) {
		err = io.EOF
		return
	}
	e = m.events[m.i]
	m.i++
	return
}

func newMockReader() *mockReader {
	return &mockReader{
		events: []mockEvent{
			{
				timestamp: time.Date(2020, 03, 12, 22, 10, 38, 0, time.UTC),
				username:  "kyle123",
				operation: "upload",
				size:      10,
			},
			{
				timestamp: time.Date(2020, 03, 12, 22, 10, 38, 0, time.UTC),
				username:  "kyle123",
				operation: "download",
				size:      20,
			},
			{
				timestamp: time.Date(2020, 04, 12, 22, 10, 38, 0, time.UTC),
				username:  "dex456",
				operation: "upload",
				size:      66,
			},
			{
				timestamp: time.Date(2020, 05, 12, 22, 10, 38, 0, time.UTC),
				username:  "dex456",
				operation: "download",
				size:      1,
			},
			{
				timestamp: time.Date(2020, 05, 13, 22, 10, 38, 0, time.UTC),
				username:  "kait789",
				operation: "download",
				size:      1024,
			},
		},
	}
}

func Test_defaultFinder_Find(t *testing.T) {
	t.Run("returns all with empty query", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, events, err := f.Find()
		assert.NoError(t, err)
		/// This is expected to return 2 because countConcern defaults to operation
		/// This should be 2 in the case where countConcern is user
		assert.Equal(t, len(r.events), count)
		assert.Equal(t, []string{
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 upload 10",
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 download 20",
			"Sun Apr 12 22:10:38 UTC 2020 dex456 upload 66",
			"Tue May 12 22:10:38 UTC 2020 dex456 download 1",
			"Wed May 13 22:10:38 UTC 2020 kait789 download 1024",
		}, events)

		fmt.Println(events)
	})

	t.Run("can count unique operations", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, _, err := f.Find(
			WithCountConcern(Operation),
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("can count users", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, _, err := f.Find(
			WithCountConcern(User),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("can match username", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, events, err := f.Find(
			WhereUsernameEquals("kyle123"),
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Equal(t, []string{
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 upload 10",
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 download 20",
		}, events)
	})

	t.Run("can match time range", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, events, err := f.Find(
			WhereTimestampIsBetween(
				time.Date(2020, 01, 01, 00, 00, 00, 0, time.UTC),
				time.Date(2020, 03, 22, 00, 00, 00, 0, time.UTC),
			),
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Equal(t, []string{
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 upload 10",
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 download 20",
		}, events)
	})

	t.Run("can match operation", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, events, err := f.Find(
			WhereOperationEquals("upload"),
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Equal(t, []string{
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 upload 10",
			"Sun Apr 12 22:10:38 UTC 2020 dex456 upload 66",
		}, events)
	})

	t.Run("can match size", func(t *testing.T) {
		r := newMockReader()
		f := NewFinder(r)
		count, events, err := f.Find(
			WhereSizeGreaterThanOrEqual(10),
			WhereSizeLessThanOrEqual(100),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
		assert.Equal(t, []string{
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 upload 10",
			"Thu Mar 12 22:10:38 UTC 2020 kyle123 download 20",
			"Sun Apr 12 22:10:38 UTC 2020 dex456 upload 66",
		}, events)
	})

	//TODO: More tests to prove query combinations

	//TODO: Test error conditions
}
