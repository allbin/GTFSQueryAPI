package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/v2/storage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type API interface {
	Close(ctx context.Context)
}

type apiRouter struct {
	conn *pgx.Conn
}

func (r *apiRouter) Close(ctx context.Context) {
	r.conn.Close(ctx)
}

func NewRouter(ctx context.Context, r *mux.Router, c config.DatabaseConfiguration) (API, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	queries := storage.New(conn)

	r.HandleFunc("/stops/{lon}/{lat}", getStopsHandler(queries)).Methods("GET")

	return &apiRouter{
		conn: conn,
	}, nil
}

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
