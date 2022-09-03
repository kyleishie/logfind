package logfind

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	errMock = Error("mock")
)

func TestNewFindOptions(t *testing.T) {
	t.Run("default cc", func(t *testing.T) {
		opts, err := newFindOptions()
		assert.NoError(t, err)
		assert.Equal(t, Event, opts.cc)
	})

	t.Run("no unexpected side effects", func(t *testing.T) {
		opts, err := newFindOptions()
		assert.NoError(t, err)
		assert.Equal(t, Event, opts.cc)
		assert.Nil(t, opts.username)
		assert.Nil(t, opts.minTime)
		assert.Nil(t, opts.maxTime)
		assert.Nil(t, opts.operation)
		assert.Nil(t, opts.minSize)
		assert.Nil(t, opts.maxSize)
	})

	t.Run("fails on first error", func(t *testing.T) {
		opts, err := newFindOptions(func(opt *finderOptions) error {
			return errMock
		})
		assert.ErrorIs(t, err, errMock)
		assert.Nil(t, opts)
	})

	t.Run("calls opt funcs", func(t *testing.T) {
		var called bool
		// Flip called each time fn is called
		fn := func(opt *finderOptions) error {
			called = !called
			return nil
		}
		// Call fn three times so that if newFindOptions calls it properly called ends up as true
		_, _ = newFindOptions(fn, fn, fn)
		assert.True(t, called)
	})
}

func TestWithCountConcern(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WithCountConcern(Operation)
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Empty(t, opt.cc)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Equal(t, Operation, opt.cc)
	})

}

func TestWithUsername(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WhereUsernameEquals("giggity")
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.username)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Equal(t, "giggity", *opt.username)
	})

}

func TestWithTimeRange(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WhereTimestampIsBetween(
			time.Date(2022, 01, 01, 00, 00, 00, 0, time.UTC),
			time.Date(2022, 12, 31, 11, 59, 59, 0, time.UTC),
		)
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.minTime)
		assert.Nil(t, opt.maxTime)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2022, 01, 01, 00, 00, 00, 0, time.UTC), *opt.minTime)
		assert.Equal(t, time.Date(2022, 12, 31, 11, 59, 59, 0, time.UTC), *opt.maxTime)
	})

	t.Run("max not before min", func(t *testing.T) {
		fn := WhereTimestampIsBetween(
			time.Date(2022, 12, 31, 11, 59, 59, 0, time.UTC),
			time.Date(2022, 01, 01, 00, 00, 00, 0, time.UTC),
		)
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.minTime)
		assert.Nil(t, opt.maxTime)

		err := fn(&opt)
		assert.ErrorIs(t, err, ErrTimeRangeInvalid)
		assert.Nil(t, opt.minTime)
		assert.Nil(t, opt.maxTime)
	})

}

func TestWithOperation(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WhereOperationEquals("upload")
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.operation)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Equal(t, "upload", *opt.operation)
	})

}

func TestWithMinSize(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WhereSizeGreaterThanOrEqual(10)
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.minSize)
		assert.Nil(t, opt.maxSize)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Equal(t, 10, *opt.minSize)
		assert.Nil(t, opt.maxSize)
	})

}

func TestWithMaxSize(t *testing.T) {

	t.Run("sets opt field(s)", func(t *testing.T) {
		fn := WhereSizeLessThanOrEqual(10)
		opt := finderOptions{}
		/// Check default state because I'm paranoid
		assert.Nil(t, opt.minSize)
		assert.Nil(t, opt.maxSize)

		err := fn(&opt)
		assert.NoError(t, err)
		assert.Nil(t, opt.minSize)
		assert.Equal(t, 10, *opt.maxSize)
	})

}
