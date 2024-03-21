package api

import (
	"context"
	"fmt"

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
	r.HandleFunc("/departures/{stop_id}", getDeparturesHandler(queries)).Methods("GET")

	return &apiRouter{
		conn: conn,
	}, nil
}
