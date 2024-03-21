package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/allbin/gtfsQueryGoApi/v2/storage"
	"github.com/gorilla/mux"
)

func getStopHandler(queries *storage.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["stop_id"]


		if id == "" {
			http.Error(w, "missing stop_id", http.StatusBadRequest)
			return
		}

		stop, err := queries.GetStop(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get stop %s: %v", id, err), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(stop)
	}
}

func getStopsHandler(queries *storage.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// get longitude and latitude from query parameters
		lon := r.URL.Query().Get("lon")
		lat := r.URL.Query().Get("lat")

		var arg storage.GetStopsNearbyParams
		if lon == "" || lat == "" {
			http.Error(w, "missing longitude or latitude", http.StatusBadRequest)
			return
		}

		arg.Lon, err = strconv.ParseFloat(lon, 64)
		if err != nil {
			http.Error(w, "invalid longitude", http.StatusBadRequest)
			return
		}

		arg.Lat, err = strconv.ParseFloat(lat, 64)
		if err != nil {
			http.Error(w, "invalid latitude", http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("radius") != "" {
			arg.Radius, err = strconv.ParseFloat(r.URL.Query().Get("radius"), 64)
		} else {
			arg.Radius = 1000
		}

		if r.URL.Query().Get("limit") != "" {
			arg.Lim, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
		} else {
			arg.Lim = 10
		}

		stops, err := queries.GetStopsNearby(r.Context(), arg)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get stops: %v", err), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(stops)
	}
}
