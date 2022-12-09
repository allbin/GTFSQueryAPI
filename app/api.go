package app

import (
	"encoding/json"
	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/direction"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/allbin/gtfsQueryGoApi/stop_departures"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	geo "github.com/martinlindhe/google-geolocate"
	"log"
	"net/http"
	"os"
	"strings"
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
	err := repo.Connect(conf.Database)
	if err != nil {
		panic(err)
	}
	log.Print("Server up and running...")
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/departures/place", placeHandler).Methods("GET")
	r.HandleFunc("/departures/stop", stopDeparturesHandler).Methods("GET")

	corsAllowedOrigins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsAllowedOrigins)(r)))
}

func placeHandler(w http.ResponseWriter, r *http.Request) {
	direction.PlaceHandler(repo, w, r, conf.Default, geoClient)
}

func stopDeparturesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	apiKey := os.Getenv("GTFS_QUERY_API_KEY")
	if apiKey != "" {
		k := q.Get("k")
		if k != apiKey {
			http.Error(w, "Missing k parameter(API KEY)", http.StatusUnauthorized)
			return
		}
	}

	stopId := q.Get("id")

	departureRows, err := stop_departures.GetStopDepartures(repo, strings.Split(stopId, ","))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_ = json.NewEncoder(w).Encode(departureRows)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//func corsMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.Method == http.MethodOptions {
//			w.Header().Add("Access-Control-Allow-Origin", "*")
//		}
//		w.WriteHeader(204)
//	})
//}
