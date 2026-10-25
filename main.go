package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/kedarui/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = ":8080"

func serveIndex() http.Handler {
	handler := http.FileServer(http.Dir("."))
	return http.StripPrefix("/app", handler)
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}
	dbQueries := database.New(db)

	tokenSecret := os.Getenv("TOKEN_SECRET")

	apiCfg := &apiConfig{
		database: dbQueries,
		platform: os.Getenv("PLATFORM"),
		tokenSecret: tokenSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serveIndex()))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", apiCfg.handlerGetChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	fmt.Println("Starting Chirpy server...")
	s := http.Server{
		Addr:    port,
		Handler: mux,
	}
	s.ListenAndServe()
}
