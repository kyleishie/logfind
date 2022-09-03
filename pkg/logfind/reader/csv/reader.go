package csv

import (
	"encoding/csv"
	lfReader "github.com/kyleishie/logfind/pkg/logfind/reader"
	"io"
	"strconv"
	"time"
)

type reader struct {
	csvReader  *csv.Reader
	pastHeader bool
}

func NewReader(r io.Reader) lfReader.Reader {
	return &reader{
		csvReader: csv.NewReader(r),
	}
}

// Read reads one event from r.
// If the record has an unexpected number of fields,
// Read returns the record along with the error ErrFieldCount.
// Except for that case, Read always returns either a non-nil
// record or a non-nil error, but not both.
// If there is no data left to be read, Read returns nil, io.EOF.
func (r *reader) Read() (e lfReader.Event, err error) {
	record, err := r.csvReader.Read()
	if !r.pastHeader && len(record) == fieldCount {
		if record[indexTimestamp] == "timestamp" &&
			record[indexUsername] == "username" &&
			record[indexOperation] == "operation" &&
			record[indexSize] == "size" {
			// Skip header
			record, err = r.csvReader.Read()
		}
	}
	e = event(record)
	return
}

const fieldCount = 4
const (
	indexTimestamp = iota
	indexUsername
	indexOperation
	indexSize
)

// event is the concrete implementation of reader.Event
// Note: Even though the fields below are based on slice index we are safe because the csv.Reader detects and errors
// on unexpected number of fields.
type event []string

func (e event) Timestamp() (timestamp time.Time, err error) {
	if len(e) != fieldCount {
		err = csv.ErrFieldCount
		return
	}
	//TODO: Make format configurable
	return time.Parse(time.UnixDate, e[indexTimestamp])
}

func (e event) Username() (username string, err error) {
	if len(e) != fieldCount {
		err = csv.ErrFieldCount
		return
	}
	return e[indexUsername], nil
}

func (e event) Operation() (op string, err error) {
	if len(e) != fieldCount {
		err = csv.ErrFieldCount
		return
	}
	return e[indexOperation], nil
}

func (e event) Size() (size int, err error) {
	if len(e) != fieldCount {
		err = csv.ErrFieldCount
		return
	}
	sizeStr := e[indexSize]
	return strconv.Atoi(sizeStr)
}
