package main

import (
	"log"
	"net/http"
)

func marshalError(w http.ResponseWriter, err error) {
	log.Printf("Error marshalling json: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func respondWithError(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if !(len(data) > 0) {
		data = []byte("Something went wrong\n")
	}
	w.Write(data)
}

func respondWithOK(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data = append(data, byte('\n'))
	w.Write(data)
}

func respondWithCreated(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data = append(data, byte('\n'))
	w.Write(data)
}

func respondWithNoContent(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	data = append(data, byte('\n'))
	w.Write(data)
}
