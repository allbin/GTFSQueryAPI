package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
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
		for _, stopDeparture := range stopDepartures {
			arv, err := gtfsTime(stopDeparture.Date, stopDeparture.Arrival)
			if err != nil {
				log.Printf("unable to parse arrival time: %v", err)
				continue
			}

			dep, err := gtfsTime(stopDeparture.Date, stopDeparture.Departure)
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
