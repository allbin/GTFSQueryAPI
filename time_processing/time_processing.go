package time_processing

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func AddTime(t time.Time, tstring string) time.Time {
	parts := strings.Split(tstring, ":")
	var intParts []time.Duration
	for i := range parts {
		val, err := strconv.Atoi(parts[i])
		if err != nil {
			panic(err)
		}
		intParts = append(intParts, time.Duration(val))
	}
	return t.Add(time.Hour * intParts[0]).Add(time.Minute * intParts[1]).Add(time.Second * intParts[2])
}

func GetTimeDifference(location *time.Location, location2 *time.Location) int {
	now := time.Now().In(location)
	return now.Hour() - time.Now().In(location2).Hour()
}

var location *time.Location

func GetLocation() *time.Location {
	if location == nil {
		location, _ = time.LoadLocation("Europe/Stockholm")
	}
	return location
}

func FromDateAndTime(d string, t string) (time.Time, error) {
	calculated, err := time.ParseInLocation("2006-01-02", d, GetLocation())
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date: %w", err)
	}

	timeParts := strings.Split(t, ":")
	if len(timeParts) != 3 {
		return time.Time{}, fmt.Errorf("unable to parse time from '%s'", t)
	}

	hour, err := strconv.ParseInt(timeParts[0], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse hour: %w", err)
	}
	minute, err := strconv.ParseInt(timeParts[1], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse minute: %w", err)
	}
	second, err := strconv.ParseInt(timeParts[2], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse second: %w", err)
	}

	return calculated.Add(
		time.Duration(hour)*time.Hour +
			time.Duration(minute)*time.Minute +
			time.Duration(second)*time.Second,
	).UTC(), nil
}
