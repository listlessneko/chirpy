package main

import (
	"encoding/json"
	"net/http"
)

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
	text := redactProfanity(c.Body)
	resp := cleaned_chirp{CleanedBody: text}
	data, err := json.Marshal(&resp)
	if err != nil {
		marshalError(w, err)
		return
	}
	respondWithOK(w, data)
}
