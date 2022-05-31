package app

import (
	"github.com/allbin/gtfsQueryGoApi/config"
	"github.com/allbin/gtfsQueryGoApi/direction"
	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/gorilla/mux"
  "github.com/gorilla/handlers"
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
    w.WriteHeader(204);
  })
}
