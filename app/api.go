package app

import (
	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/direction"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/gorilla/mux"
	geo "github.com/martinlindhe/google-geolocate"
	"log"
	"net/http"
	"os"
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
	log.Fatal(http.ListenAndServe(":8080", r))
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
