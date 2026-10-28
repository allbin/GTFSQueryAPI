package health

import (
	"encoding/json"
	"net/http"
)

type status struct {
	Status string `json:"status"`
}

// Handler answers the Kubernetes readiness and liveness probes. It does not
// touch the database: the probes only need to know that the process serves HTTP.
func Handler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(status{Status: "ok"})
}
