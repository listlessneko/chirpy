package main

import (
	"fmt"
	"net/http"
)

const port = ":8080"

func serveIndex() http.Handler {
	handler := http.FileServer(http.Dir("."))
	return http.StripPrefix("/app", handler)
}

func main() {
	apiCfg := &apiConfig{}

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
