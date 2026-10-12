package main

import (
	"fmt"
	"net/http"
)

const filepathRoot = "."
const filepathApp = "/app/"
const port = ":8080"

func serveIndex() http.Handler {
	handler := http.FileServer(http.Dir(filepathRoot))
	return http.StripPrefix("/app", handler)
}

func main() {
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serveIndex()))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

	fmt.Println("Starting Chirpy server...")
	s := http.Server{
		Addr:    port,
		Handler: mux,
	}
	s.ListenAndServe()
}
