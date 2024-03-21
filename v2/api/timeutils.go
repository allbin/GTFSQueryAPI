package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func dateParts(date string) (int, int, int, error) {
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid date format")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid year")
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid month")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid day")
	}

	return year, month, day, nil
}

func timeParts(time string) (int, int, int, error) {
	parts := strings.Split(time, ":")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid hour")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid minute")
	}

	second, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid second")
	}

	return hour, minute, second, nil
}

func gtfsTime(date string, tm string) (time.Time, error) {
	year, month, day, err := dateParts(date)
	if err != nil {
		return time.Time{}, err
	}

	hour, minute, second, err := timeParts(tm)
	if err != nil {
		return time.Time{}, err
	}

	t := time.Date(year, time.Month(month), day, 0, minute, second, 0, time.Local)
	return t.Add(time.Duration(hour) * time.Hour), nil
}
