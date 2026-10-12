package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	status := fmt.Sprintf("Status: %d OK\n", http.StatusOK)
	w.Write([]byte(status))
}
