package logfind

import "time"

type finderOptions struct {
	cc       CountConcern
	username *string

	minTime *time.Time
	maxTime *time.Time

	operation *string

	minSize *int
	maxSize *int
}

func newFindOptions(opts ...FinderOptionFunc) (opt *finderOptions, err error) {
	opt = &finderOptions{
		cc: Event,
	}

	for _, optionFunc := range opts {
		err = optionFunc(opt)
		if err != nil {
			opt = nil
			return
		}
	}

	return
}

type FinderOptionFunc func(*finderOptions) error

// CountConcern is your entrypoint to customize how events are counted.
type CountConcern string

const (
	// Event - When event is used, events are counted regardless of values. This is the default.
	Event = CountConcern("event")
	// Operation - When operation is used, events are counted by unique operation values.
	Operation = CountConcern("operation")

	// User - When user is used, events are counted by unique username values.
	User = CountConcern("user")
)

// WithCountConcern customizes how the Finder counts events. See the constant CountConcerns for details.
func WithCountConcern(concern CountConcern) FinderOptionFunc {
	return func(opt *finderOptions) error {
		opt.cc = concern
		return nil
	}
}

// WhereUsernameEquals adds the requirement that matching log events must contain the given username.
func WhereUsernameEquals(username string) FinderOptionFunc {
	return func(opt *finderOptions) error {
		opt.username = &username
		return nil
	}
}

// WhereTimestampIsBetween adds the requirement that matching log events timestamp value is between the given time range.
func WhereTimestampIsBetween(min, max time.Time) FinderOptionFunc {
	return func(opt *finderOptions) error {
		if max.Before(min) {
			return ErrTimeRangeInvalid
		}

		opt.minTime = &min
		opt.maxTime = &max
		return nil
	}
}

// WhereOperationEquals adds the requirement that matching log events operation value is equal to the given operation string.
func WhereOperationEquals(operation string) FinderOptionFunc {
	return func(opt *finderOptions) error {
		opt.operation = &operation
		return nil
	}
}

// WhereSizeGreaterThanOrEqual adds the requirement that matching log events size value is greater than or equal to min.
//
// Note: Size is represented in kB.
func WhereSizeGreaterThanOrEqual(min int) FinderOptionFunc {
	return func(opt *finderOptions) error {
		opt.minSize = &min
		return nil
	}
}

// WhereSizeLessThanOrEqual adds the requirement that matching log events size value is less than or equal to max.
//
// Note: Size is represented in kB.
func WhereSizeLessThanOrEqual(max int) FinderOptionFunc {
	return func(opt *finderOptions) error {
		opt.maxSize = &max
		return nil
	}
}
