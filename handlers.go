package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Body string `json:"body"`
	}

	type SuccessReturnVals struct {
		SuccessMessage 	string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	const chirpCharLimit = 140
	if len(params.Body) > chirpCharLimit {	
		JSONErrorResponse(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	} 

	profane_replacement := "****"
	profane_list := []string{"kerfuffle", "sharbert", "fornax"}
	bodyMsg := strings.Split(params.Body, " ")
	for i, word := range bodyMsg {
		for _, invalid_word := range profane_list{
			if invalid_word == strings.ToLower(word) {
				bodyMsg[i] = profane_replacement
			}
		}
	}
	cleanedMessage := strings.Join(bodyMsg, " ")
	respBody := SuccessReturnVals{
		SuccessMessage: cleanedMessage,
	}
	respondWithJSON(w, http.StatusOK, respBody)
	
}