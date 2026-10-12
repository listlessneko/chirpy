package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const filepathRoot = "."
const filepathApp = "/app/"
const port = ":8080"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handlerHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func serveIndex() http.Handler {
	handler := http.FileServer(http.Dir(filepathRoot))
	return http.StripPrefix("/app", handler)
}

func main() {
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serveIndex()))
	mux.HandleFunc("/healthz", handlerHealthz)
	s := http.Server{
		Addr: port,
		Handler: mux,
	}
	s.ListenAndServe()
}
