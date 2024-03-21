package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/allbin/gtfsQueryGoApi/v2/storage"
	"github.com/gorilla/mux"
)

func getStopsHandler(queries *storage.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		params := mux.Vars(r)

		if params["lon"] == "" || params["lat"] == "" {
			http.Error(w, "missing longitude or latitude", http.StatusBadRequest)
			return
		}

		var arg storage.GetStopsNearbyParams

		arg.Lon, err = strconv.ParseFloat(params["lon"], 64)
		if err != nil {
			http.Error(w, "invalid longitude", http.StatusBadRequest)
			return
		}

		arg.Lat, err = strconv.ParseFloat(params["lat"], 64)
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
