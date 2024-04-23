package api

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/v2/storage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API interface {
	Close(ctx context.Context)
}

type apiRouter struct {
	conn *pgxpool.Pool
}

func (r *apiRouter) Close(ctx context.Context) {
	r.conn.Close()
}

func NewRouter(ctx context.Context, r *mux.Router, c config.DatabaseConfiguration) (API, error) {
	conn, err := pgxpool.New(ctx, dbString(c))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	queries := storage.New(conn)

	r.HandleFunc("/stops", getStopsHandler(queries)).Methods("GET")
	r.HandleFunc("/stops/{stop_id}", getStopHandler(queries)).Methods("GET")
	r.HandleFunc("/stops/{stop_id}/departures", getDeparturesHandler(queries)).Methods("GET")
	r.HandleFunc("/departures", getDeparturesForStopsHandler(queries)).Methods("GET")

	return &apiRouter{
		conn: conn,
	}, nil
}

func dbString(c config.DatabaseConfiguration) string {
	passwordArg := ""
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = c.Password
	}
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = c.Host
	}
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = c.Port
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = c.User
	}
	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = c.Database
	}
	if len(pass) > 0 {
		passwordArg = "password=" + pass
	}
	db_string := fmt.Sprintf("host=%s port=%d user=%s %s dbname=%s sslmode=disable",
		host, port, user, passwordArg, db)

	return db_string
}
