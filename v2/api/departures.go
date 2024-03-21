package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/allbin/gtfsQueryGoApi/v2/storage"
	"github.com/gorilla/mux"
)

type Departure struct {
	storage.GetStopDeparturesRow
	NextArrival   time.Time `json:"next_arrival"`
	NextDeparture time.Time `json:"next_departure"`
}

func getDeparturesHandler(queries *storage.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		params := mux.Vars(r)

		if params["stop_id"] == "" {
			http.Error(w, "missing stop_id", http.StatusBadRequest)
			return
		}

		var arg storage.GetStopDeparturesParams

		arg.StopID = params["stop_id"]

		if r.URL.Query().Get("limit") != "" {
			arg.Lim, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
		} else {
			arg.Lim = 1000
		}

		stopDepartures, err := queries.GetStopDepartures(r.Context(), arg)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get departures: %v", err), http.StatusInternalServerError)
		}

		var departures []Departure
		now := time.Now().Local()

		for _, stopDeparture := range stopDepartures {
			arv, err := GtfsTime(stopDeparture.Arrival).Time(now)
			if err != nil {
				log.Printf("unable to parse arrival time: %v", err)
				continue
			}

			dep, err := GtfsTime(stopDeparture.Departure).Time(now)
			if err != nil {
				log.Printf("unable to parse departure time: %v", err)
				continue
			}

			departures = append(departures, Departure{
				GetStopDeparturesRow: stopDeparture,
				NextArrival:          arv,
				NextDeparture:        dep,
			})
		}
		slices.SortFunc(departures, func(i, j Departure) int {
			if i.NextDeparture.Before(j.NextDeparture) {
				return -1
			} else if i.NextDeparture.After(j.NextDeparture) {
				return 1
			}
			return 0
		})

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(departures)

	}
}

type GtfsTime string

func (t GtfsTime) Time(offset time.Time) (time.Time, error) {
	parts := strings.Split(string(t), ":")

	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute")
	}

	second, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid second")
	}

	year, month, day := offset.Date()
	timestamp := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	timestamp = timestamp.Add(time.Duration(hour) * time.Hour)
	timestamp = timestamp.Add(time.Duration(minute) * time.Minute)
	timestamp = timestamp.Add(time.Duration(second) * time.Second)

	if timestamp.Before(offset) {
		timestamp = timestamp.AddDate(0, 0, 1)
	}

	return timestamp, nil
}
