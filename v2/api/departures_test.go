package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGTFSTime(t *testing.T) {
	year, month, day := time.Now().Local().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	t.Run("Simple", func(t *testing.T) {
		input := GtfsTime("12:34:56")
		ts, err := input.Time(year, int(month), day)

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
		ts, err := input.Time(year, int(month), day)

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
		ts, err := input.Time(year, int(month), day)

		expected := today.AddDate(0, 0, 1)
		expected = expected.Add(2 * time.Hour)
		expected = expected.Add(34 * time.Minute)
		expected = expected.Add(56 * time.Second)

		assert.NoError(t, err)
		assert.Equal(t, expected, ts)
	})
}
