package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGTFSTime(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ts, err := gtfsTime("2024-01-02", "12:34:56")

		assert.NoError(t, err)
		assert.Equal(t, 2024, ts.Year())
		assert.Equal(t, 01, int(ts.Month()))
		assert.Equal(t, 02, ts.Day())
		assert.Equal(t, 12, ts.Hour())
		assert.Equal(t, 34, ts.Minute())
		assert.Equal(t, 56, ts.Second())
	})

	t.Run("Afternoon", func(t *testing.T) {
		ts, err := gtfsTime("2024-01-02", "22:34:56")

		assert.NoError(t, err)
		assert.Equal(t, 2024, ts.Year())
		assert.Equal(t, 01, int(ts.Month()))
		assert.Equal(t, 02, ts.Day())
		assert.Equal(t, 22, ts.Hour())
		assert.Equal(t, 34, ts.Minute())
		assert.Equal(t, 56, ts.Second())
	})

	t.Run("Next Day", func(t *testing.T) {
		ts, err := gtfsTime("2024-01-02", "26:34:56")

		assert.NoError(t, err)
		assert.Equal(t, 2024, ts.Year())
		assert.Equal(t, 01, int(ts.Month()))
		assert.Equal(t, 03, ts.Day())
		assert.Equal(t, 02, ts.Hour())
		assert.Equal(t, 34, ts.Minute())
		assert.Equal(t, 56, ts.Second())
	})
}
