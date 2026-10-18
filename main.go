package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/joho/godotenv"
)

const port = ":8080"

func serveIndex() http.Handler {
	handler := http.FileServer(http.Dir("."))
	return http.StripPrefix("/app", handler)
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{database: dbQueries}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serveIndex()))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	fmt.Println("Starting Chirpy server...")
	s := http.Server{
		Addr:    port,
		Handler: mux,
	}
	s.ListenAndServe()
}
