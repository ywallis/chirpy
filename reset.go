package main

import "net/http"

func (a *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	a.fileServerHits.Store(0)
	w.Write([]byte("Hit counter reset"))

}
