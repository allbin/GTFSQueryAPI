package stop_departures

import (
	"fmt"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/allbin/gtfsQueryGoApi/time_processing"
	"time"
)

type Departure struct {
	ArrivalTime   time.Time `json:"arrival_time"`
	DepartureTime time.Time `json:"departure_time"`
	Headsign      string    `json:"headsign"`
	ShortName     string    `json:"short_name"`
	LongName      string    `json:"long_name"`
}

type StopDeparture struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Lat        float64     `json:"lat"`
	Lon        float64     `json:"lon"`
	Departures []Departure `json:"departures"`
}

type StopDepartures = map[string]StopDeparture

func GetStopDepartures(repo *query.Repository, stops []string) (StopDepartures, error) {
	//startTime := time.Now()
	//defer func() {
	//	duration := time.Since(startTime)
	//	log.Printf("GetStopDepartures: %dms", duration.Milliseconds())
	//}()

	rows, err := repo.GetDepartures(stops)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var departureRows []query.DepartureRow

	for rows.Next() {
		var departure query.DepartureRow
		err := rows.Scan(
			&departure.Id,
			&departure.ArrivalTime,
			&departure.DepartureTime,
			&departure.Name,
			&departure.Lat,
			&departure.Lon,
			&departure.Headsign,
			&departure.ShortName,
			&departure.LongName,
			&departure.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan departure: %w", err.Error())
		}

		departureRows = append(departureRows, departure)
	}

	stopDepartures := StopDepartures{}

	for _, row := range departureRows {
		stopDeparture, ok := stopDepartures[row.Id]
		if !ok {
			// StopDeparture missing, create new
			stopDeparture = StopDeparture{
				Id:         row.Id,
				Name:       row.Name,
				Lat:        row.Lat,
				Lon:        row.Lon,
				Departures: []Departure{},
			}
		}

		// Parse arrival and departure time
		arrivalTime, err := time_processing.FromDateAndTime(row.Date, row.ArrivalTime)
		if err != nil {
			return nil, fmt.Errorf("error in arrival time: %w", err)
		}
		departureTime, err := time_processing.FromDateAndTime(row.Date, row.DepartureTime)
		if err != nil {
			return nil, fmt.Errorf("error in departure time: %w", err)
		}

		// Create departure and add to stopDeparture
		stopDeparture.Departures = append(stopDeparture.Departures, Departure{
			ArrivalTime:   arrivalTime,
			DepartureTime: departureTime,
			Headsign:      row.Headsign,
			ShortName:     row.ShortName,
			LongName:      row.LongName,
		})

		// Save stopDeparture back into lookup
		stopDepartures[row.Id] = stopDeparture
	}

	return stopDepartures, nil
}
