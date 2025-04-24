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
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening DB")
	}
	dbQueries := database.New(db)
	const filePathRoot string = "."
	const port string = "8080"
	apiCfg := apiConfig{db: dbQueries}
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
	server.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)
	server.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	server.Handle("GET /api/hello", middlewareLogger(middlewareAuth(http.HandlerFunc(helloHandler))))
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Error while starting server: %s\n", err)
	}
}
