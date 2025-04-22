package main

import (
	"fmt"
	"net/http"
)

const metricsHtml string = "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"




func (a *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	hits := a.fileServerHits.Load()
	w.Write(fmt.Appendf(nil, metricsHtml, hits))
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		a.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
