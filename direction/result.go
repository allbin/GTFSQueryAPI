package direction

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"

	"github.com/allbin/gtfsQueryGoApi/query"
	"github.com/allbin/gtfsQueryGoApi/time_processing"
	"github.com/cornelk/hashmap"
)

func GetResult(r *query.Repository, la float64, lo float64, radius int, maxDepartures int, maxStops int) []Result {
	rows, err := r.GetStops(strconv.FormatFloat(la, 'f', -1, 64),
		strconv.FormatFloat(lo, 'f', -1, 64), strconv.Itoa(radius), strconv.Itoa(maxStops))
	if err != nil {
		panic(err)
	}
	return groupAndSortRows(rows, maxStops, maxDepartures)
}

func groupAndSortRows(rows *sql.Rows, maxStops int, maxDepartures int) []Result {
	resultMap := hashmap.NewSized[string, Result](uintptr(maxStops * 2))
	keysOrder := []string{}

	for rows.Next() {
		var row row
		if err := rows.Scan(&row.id, &row.arrivalTime, &row.departureTime, &row.name, &row.lat, &row.lon, &row.headsign, &row.short_name, &row.long_name, &row.date, &row.dateString); err != nil {
			log.Fatal(err)
		}
		loc_name := "Europe/Stockholm"
		loc, err := time.LoadLocation(loc_name)
		if err != nil {
			panic(fmt.Sprintf("Problem loading location %s", loc_name))
		}
		timeDiff := time_processing.GetTimeDifference(loc, time.UTC)
		now := time.Now().In(time.UTC)
		date, _ := time.Parse(time.RFC3339, row.date)
		dep := time_processing.AddTime(date, row.departureTime).Add(time.Hour * time.Duration(-timeDiff))
		arr := time_processing.AddTime(date, row.arrivalTime).Add(time.Hour * time.Duration(-timeDiff))

		if dep.After(now) {
			value, exist := resultMap.Get(row.id)
			if exist == true {
				if len(value.Departures) < maxDepartures {
					value.Departures = append(value.Departures, Departure{dep.Format("15:04:05"), arr.Format("15:04:05"), dep.Format("2006-01-02T15:04:05-07:00"), Trip{row.headsign, row.short_name, row.long_name}})
					resultMap.Set(row.id, value)
				}

			} else {
				resultMap.Insert(row.id, rowToresult(row, arr, dep))
				keysOrder = append(keysOrder, row.id)
			}
		}
	}
	var r []Result
	for _, key := range keysOrder {
		value, exist := resultMap.Get(key)
		if exist == true {
			r = append(r, value)
		}
	}
	return r
}
func rowToresult(r row, arr time.Time, dep time.Time) Result {
	return Result{
		Stop{r.id, []string{r.lat, r.lon}, r.name},
		[]Departure{{dep.Format("15:04:05"), arr.Format("15:04:05"), dep.Format("2006-01-02T15:04:05-07:00"),
			Trip{r.headsign, r.short_name, r.long_name}}}}
}
