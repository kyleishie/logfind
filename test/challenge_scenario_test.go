package test

import (
	"github.com/kyleishie/logfind/pkg/logfind"
	"github.com/kyleishie/logfind/pkg/logfind/reader/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_Challenge_Scenarios(t *testing.T) {

	t.Run("How many users accessed the server?", func(t *testing.T) {
		csvFile, err := os.Open("../server_log.csv")
		if err != nil {
			t.Error(err)
		}
		r := csv.NewReader(csvFile)
		f := logfind.NewFinder(r)
		count, _, err := f.Find(
			logfind.WithCountConcern(logfind.User),
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, count)
	})

	t.Run("How many uploads were larger than 50kB?", func(t *testing.T) {
		csvFile, err := os.Open("../server_log.csv")
		if err != nil {
			t.Error(err)
		}
		r := csv.NewReader(csvFile)
		f := logfind.NewFinder(r)
		count, _, err := f.Find(
			logfind.WhereOperationEquals("upload"),
			logfind.WhereSizeGreaterThanOrEqual(50),
		)
		assert.NoError(t, err)
		assert.Equal(t, 144, count)
	})

	t.Run("How many times did jeff22 upload to the server on April 15th, 2020?", func(t *testing.T) {
		csvFile, err := os.Open("../server_log.csv")
		if err != nil {
			t.Error(err)
		}
		r := csv.NewReader(csvFile)
		f := logfind.NewFinder(r)
		count, _, err := f.Find(
			logfind.WhereUsernameEquals("jeff22"),
			logfind.WhereOperationEquals("upload"),
			logfind.WhereTimestampIsBetween(
				time.Date(2020, 04, 15, 00, 00, 00, 0, time.UTC),
				time.Date(2020, 04, 16, 00, 00, 00, 0, time.UTC),
			),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

}
