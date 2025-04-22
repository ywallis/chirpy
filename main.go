package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {

	const filePathRoot string = "."
	const port string = "8080"
	apiCfg := apiConfig{}
	server := http.NewServeMux()

	s := http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	fileHandler := http.FileServer(http.Dir(filePathRoot))
	server.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))
	server.HandleFunc("GET /api/healthz", readinessHandler)
	server.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	server.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)
	server.Handle("GET /api/hello", middlewareLogger(middlewareAuth(http.HandlerFunc(helloHandler))))
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Error while starting server: %s\n", err)
	}
}
