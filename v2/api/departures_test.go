package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGTFSTime(t *testing.T) {
	year, month, day := time.Now().Local().Date()
	today := time.Date(year, month, day, 12, 0, 0, 0, time.Local)

	t.Run("Simple", func(t *testing.T) {
		input := GtfsTime("12:34:56")
		ts, err := input.Time(today)

		assert.NoError(t, err)
		assert.Equal(t, year, ts.Year())
		assert.Equal(t, month, ts.Month())
		assert.Equal(t, day, ts.Day())
		assert.Equal(t, 12, ts.Hour())
		assert.Equal(t, 34, ts.Minute())
		assert.Equal(t, 56, ts.Second())
	})

	t.Run("Afternoon", func(t *testing.T) {
		input := GtfsTime("22:34:56")
		ts, err := input.Time(today)

		assert.NoError(t, err)
		assert.Equal(t, year, ts.Year())
		assert.Equal(t, month, ts.Month())
		assert.Equal(t, day, ts.Day())
		assert.Equal(t, 22, ts.Hour())
		assert.Equal(t, 34, ts.Minute())
		assert.Equal(t, 56, ts.Second())
	})

	t.Run("Next Day", func(t *testing.T) {
		input := GtfsTime("26:34:56")
		ts, err := input.Time(today)

		expected := time.Date(year, month, day, 2, 34, 56, 0, time.Local)
		expected = expected.AddDate(0, 0, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, ts)
	})

	t.Run("Next Day Wrapped", func(t *testing.T) {
		input := GtfsTime("11:34:56")
		ts, err := input.Time(today)

		expected := time.Date(year, month, day, 11, 34, 56, 0, time.Local)
		expected = expected.AddDate(0, 0, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, ts)
	})
}
