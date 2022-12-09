package query

import (
	"database/sql"
	"github.com/lib/pq"
)

type DepartureRow struct {
	Id            string  `json:"id,omitempty"`
	ArrivalTime   string  `json:"arrival_time,omitempty"`
	DepartureTime string  `json:"departure_time,omitempty"`
	Name          string  `json:"name,omitempty"`
	Lat           float64 `json:"lat,omitempty"`
	Lon           float64 `json:"lon,omitempty"`
	Headsign      string  `json:"headsign,omitempty"`
	ShortName     string  `json:"short_name,omitempty"`
	LongName      string  `json:"long_name"`
	Date          string  `json:"date,omitempty"`
}

func (r *Repository) GetDepartures(ids []string) (*sql.Rows, error) {
	return r.Db.Query(`
		SELECT s.stop_id                              as id,
			   arrival_time,
			   departure_time,
			   stop_name                              as name,
			   stop_lat                               as lat,
			   stop_lon                               as lon,
			   trip_headsign                          as headsign,
			   coalesce(route_short_name, '')         as short_name,
			   coalesce(route_long_name, '')          as long_name,
			   cast(date as text)                     as date
		from stop_times
				 JOIN stops s ON s.stop_id = stop_times.stop_id
				 JOIN trips t on stop_times.trip_id = t.trip_id
				 JOIN routes r ON r.route_id = t.route_id
				 JOIN calendar_dates cd on t.service_id = cd.service_id
		WHERE s.stop_id = ANY($1) AND date >= CURRENT_DATE AND date <= CURRENT_DATE + 1
		ORDER BY date, departure_time
	`, pq.Array(ids))
}
