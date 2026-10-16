package main

import (
	"encoding/json"
	// "errors"
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
	data = append(data, byte('\n'))
	w.Write(data)
}

func respondWithOK(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data = append(data, byte('\n'))
	w.Write(data)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if len(c.Body) > 140 {
		tooLongResp := errResp{Error: "Chirp is too long"}
		data, err := json.Marshal(&tooLongResp)
		if err != nil {
			marshalError(w, err)
			return
		}
		respondWithError(w, http.StatusBadRequest, data)
		return
	}
	if err != nil {
		errResp := errResp{Error: "Something went wrong"}
		data, err := json.Marshal(&errResp)
		if err != nil {
			marshalError(w, err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, data)
		return
	}
	validResp := validResp{Valid: true}
	data, err := json.Marshal(&validResp)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}
