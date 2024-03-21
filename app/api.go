package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/direction"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/allbin/gtfsQueryGoApi/v2/api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	geo "github.com/martinlindhe/google-geolocate"
)

var (
	conf      *config.Configuration
	repo      *query.Repository
	geoClient *geo.GoogleGeo
)

func init() {
	conf = new(config.Configuration)
	repo = new(query.Repository)
	geoClient = geo.NewGoogleGeo(os.Getenv("GOOGLE_GEOCODE_API_KEY"))
	err := config.Init(conf)
	if err != nil {
		panic(err)
	}
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := repo.Connect(conf.Database)
	if err != nil {
		panic(err)
	}
	log.Print("Server up and running...")
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/departures/place", placeHandler).Methods("GET")

	v2Router := r.PathPrefix("/api/v2").Subrouter()
	v2Api, err := api.NewRouter(ctx, v2Router, conf.Database)
	if err != nil {
		panic(err)
	}
	defer v2Api.Close(ctx)

	corsAllowedOrigins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsAllowedOrigins)(r)))
}

func placeHandler(w http.ResponseWriter, r *http.Request) {
	direction.PlaceHandler(repo, w, r, conf.Default, geoClient)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(204)
	})
}
