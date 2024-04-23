package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
	log.Warn("status", "status", status)
}

func httpRequestEntry(r *http.Request, status int, latency string) map[string]interface{} {
	// The log entry is in JSON format and has the following fields:
	// {
	// 	"requestMethod": string,
	// 	"requestUrl": string,
	// 	"requestSize": string,
	// 	"status": integer,
	// 	"responseSize": string,
	// 	"userAgent": string,
	// 	"remoteIp": string,
	// 	"serverIp": string,
	// 	"referer": string,
	// 	"latency": string,
	// 	"cacheLookup": boolean,
	// 	"cacheHit": boolean,
	// 	"cacheValidatedWithOriginServer": boolean,
	// 	"cacheFillBytes": string,
	// 	"protocol": string
	// }

	entry := map[string]interface{}{}
	entry["requestMethod"] = r.Method
	entry["requestUrl"] = r.URL.String()
	entry["requestSize"] = r.ContentLength
	entry["status"] = status
	entry["userAgent"] = r.UserAgent()
	entry["remoteIp"] = r.RemoteAddr
	entry["serverIp"] = r.Host
	entry["latency"] = latency
	entry["protocol"] = r.Proto

	return entry
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &ResponseWriter{ResponseWriter: w, Status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(writer, r)

		keyvals := []interface{}{
			"httpRequest", httpRequestEntry(r, writer.Status, time.Now().Sub(start).String()),
		}
		// keyvals = append(keyvals, "httpRequest", httpRequestEntry(r, writer.Status, time.Now().Sub(start).String()))
		origin := r.URL.Query().Get("origin")
		if len(origin) > 0 {
			keyvals = append(keyvals, "origin", origin)
		}

		log.Info(
			fmt.Sprintf("%s %s", r.Method, r.URL),
			keyvals...,
		)
	})

}
