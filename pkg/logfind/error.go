package logfind

type Error string

var _ error = (*Error)(nil)

func (e Error) Error() string {
	return string(e)
}

const (
	ErrTimeRangeInvalid = Error("time range invalid")
	ErrSizeRangeInvalid = Error("size range invalid")
)
