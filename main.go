package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ywallis/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db *database.Queries
	platform string
	JWTSecret string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	JWTSecret := os.Getenv("JWTSECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening DB")
	}
	dbQueries := database.New(db)
	const filePathRoot string = "."
	const port string = "8080"
	apiCfg := apiConfig{db: dbQueries, platform: platform, JWTSecret: JWTSecret}
	server := http.NewServeMux()

	s := http.Server{
		Addr:    ":" + port,
		Handler: server,
	}

	fileHandler := http.FileServer(http.Dir(filePathRoot))
	server.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))
	server.HandleFunc("GET /api/healthz", readinessHandler)
	server.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	server.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	server.HandleFunc("PUT /api/users", apiCfg.handlerUpdateCredentials)
	server.HandleFunc("POST /admin/reset", apiCfg.handlerReset) 
	server.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	server.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	server.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	server.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	server.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	server.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	server.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)
	server.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerDeleteChirp)
	server.Handle("GET /api/hello", middlewareLogger(middlewareAuth(http.HandlerFunc(helloHandler))))
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Error while starting server: %s\n", err)
	}
}
