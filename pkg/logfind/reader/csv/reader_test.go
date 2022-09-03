package csv

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	t.Run("respects given reader", func(t *testing.T) {
		input := strings.NewReader("Sun Apr 12 22:10:38 UTC 2020,sarah94,download,34")
		r := NewReader(input)
		e, err := r.Read()
		assert.Nil(t, err)
		assert.Equal(t, event([]string{"Sun Apr 12 22:10:38 UTC 2020", "sarah94", "download", "34"}), e)
	})
}

func Test_event_Timestamp(t *testing.T) {
	t.Run("can parse timestamp as Unix Date", func(t *testing.T) {
		e := event([]string{"Sun Apr 12 22:10:38 UTC 2020", "sarah94", "download", "34"})
		gotTimestamp, err := e.Timestamp()
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2020, 04, 12, 22, 10, 38, 0, time.UTC), gotTimestamp)
	})

	t.Run("fails on unexpected timestamp format", func(t *testing.T) {
		e := event([]string{"2022-01-01T00:00:00.000Z", "sarah94", "download", "34"})
		gotTimestamp, err := e.Timestamp()
		assert.Error(t, err)
		assert.Empty(t, gotTimestamp)
	})

	t.Run("respects expected field count", func(t *testing.T) {
		e := event([]string{})
		gotTimestamp, err := e.Timestamp()
		assert.ErrorIs(t, err, csv.ErrFieldCount)
		assert.Empty(t, gotTimestamp)
	})
}

func Test_event_Username(t *testing.T) {
	t.Run("can parse username", func(t *testing.T) {
		e := event([]string{"Sun Apr 12 22:10:38 UTC 2020", "sarah94", "download", "34"})
		gotUsername, err := e.Username()
		assert.NoError(t, err)
		assert.Equal(t, "sarah94", gotUsername)
	})

	t.Run("respects expected field count", func(t *testing.T) {
		e := event([]string{})
		gotUsername, err := e.Username()
		assert.ErrorIs(t, err, csv.ErrFieldCount)
		assert.Empty(t, gotUsername)
	})
}

func Test_event_Operation(t *testing.T) {
	t.Run("can parse username", func(t *testing.T) {
		e := event([]string{"Sun Apr 12 22:10:38 UTC 2020", "sarah94", "download", "34"})
		gotOperation, err := e.Operation()
		assert.NoError(t, err)
		assert.Equal(t, "download", gotOperation)
	})

	t.Run("respects expected field count", func(t *testing.T) {
		e := event([]string{})
		gotOperation, err := e.Operation()
		assert.ErrorIs(t, err, csv.ErrFieldCount)
		assert.Empty(t, gotOperation)
	})
}

func Test_event_Size(t *testing.T) {
	t.Run("can parse username", func(t *testing.T) {
		e := event([]string{"Sun Apr 12 22:10:38 UTC 2020", "sarah94", "download", "34"})
		gotSize, err := e.Size()
		assert.NoError(t, err)
		assert.Equal(t, 34, gotSize)
	})

	t.Run("respects expected field count", func(t *testing.T) {
		e := event([]string{})
		gotSize, err := e.Size()
		assert.ErrorIs(t, err, csv.ErrFieldCount)
		assert.Empty(t, gotSize)
	})
}
