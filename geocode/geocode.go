package geocode

import (
	"github.com/charmbracelet/log"
	geo "github.com/martinlindhe/google-geolocate"
)

func GetCordinatesForAddress(address string, geo *geo.GoogleGeo) (lat float64, lon float64) {
	res, err := geo.Geocode(address)
	if err != nil {
		log.Error("failed to geocode address", "address", address, "error", err)
		return 0, 0
	}
	return res.Lat, res.Lng
}
