package direction

type Result struct {
	Stop       Stop        `json:"stop"`
	Departures []Departure `json:"departures"`
}

type Stop struct {
	Id   string   `json:"id"`
	Loc  []string `json:"loc"`
	Name string   `json:"name"`
}
type Departure struct {
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Date          string `json:"date"`
	Trip          Trip   `json:"trip"`
}
type Trip struct {
	Headsign       string `json:"headsign"`
	RouteShortName string `json:"short_name"`
	RouteLongName  string `json:"long_name"`
}

type row struct {
	id            string
	arrivalTime   string
	departureTime string
	name          string
	lat           string
	lon           string
	headsign      string
	short_name    string
	long_name     string
	date          string
	dateString    string
}
