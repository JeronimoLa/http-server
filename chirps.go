package main

import (
	"encoding/json"
	"net/http"
	"strings"


	"github.com/jeronimoLa/http-server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Body 	string `json:"body"`
	UserID  uuid.UUID `json:"user_id"`
}

func NewChirpResponse(u database.Chirp) ChirpResponse {
	return ChirpResponse{
		ChirpID: 	u.ChirpID,
		CreatedAt: 	u.CreatedAt,
		UpdatedAt: 	u.UpdatedAt,
		Body: 		u.Body,	
		ID:			u.ID,
	}
}

func validateChirp(w http.ResponseWriter, ReqBody *Chirp) database.AddChirpsToUserParams { 
	const chirpCharLimit = 140
	if len(ReqBody.Body) > chirpCharLimit {	
		JSONErrorResponse(w, http.StatusBadRequest, "Chirp is too long", nil)
		return database.AddChirpsToUserParams{}
	} 

	profane_replacement := "****"
	profane_list := []string{"kerfuffle", "sharbert", "fornax"}
	bodyMsg := strings.Split(ReqBody.Body, " ")
	for i, word := range bodyMsg {
		for _, invalid_word := range profane_list{
			if invalid_word == strings.ToLower(word) {
				bodyMsg[i] = profane_replacement
			}
		}
	}
	cleanedMessage := strings.Join(bodyMsg, " ")
	params := database.AddChirpsToUserParams{
		Body: cleanedMessage,
		ID:	ReqBody.UserID,
	}
	return params
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	ReqBody := Chirp{}
	err := decoder.Decode(&ReqBody)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	params := validateChirp(w, &ReqBody)

	chirpDetails, err := cfg.db.AddChirpsToUser(r.Context(), params)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Something went wrong with updating chirp to user", err)
		return
	}
	respBody := NewChirpResponse(chirpDetails)
	respondWithJSON(w, http.StatusCreated, respBody)
}