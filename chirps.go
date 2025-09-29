package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	// "log"
	"github.com/google/uuid"
	"github.com/jeronimoLa/http-server/internal/auth"
	"github.com/jeronimoLa/http-server/internal/database"
)

type Chirp struct {
	Body 	string    `json:"body"`
}

type ChirpResponse struct {
	ID   		uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Body      	string    `json:"body"`
	UserID      uuid.UUID `json:"user_id"`
}

func NewChirpResponse(u database.Chirp) ChirpResponse {
	return ChirpResponse{
		ID:   	   u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Body:      u.Body,
		UserID:    u.UserID,
	}
}

func validateChirp(ReqBody *Chirp) (database.AddChirpsToUserParams, error) {
	const chirpCharLimit = 140
	if len(ReqBody.Body) == 0 {
		return database.AddChirpsToUserParams{}, fmt.Errorf("body key-valye is empty")
	}

	if len(ReqBody.Body) > chirpCharLimit {
		return database.AddChirpsToUserParams{}, fmt.Errorf("Chirp is too long")
	}

	profane_replacement := "****"
	profane_list := []string{"kerfuffle", "sharbert", "fornax"}
	bodyMsg := strings.Split(ReqBody.Body, " ")
	for i, word := range bodyMsg {
		for _, invalid_word := range profane_list {
			if invalid_word == strings.ToLower(word) {
				bodyMsg[i] = profane_replacement
			}
		}
	}
	cleanedMessage := strings.Join(bodyMsg, " ")
	params := database.AddChirpsToUserParams{
		Body: cleanedMessage,
	}
	return params, nil
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	ReqBody := Chirp{}
	err = decoder.Decode(&ReqBody)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't validate token", err)
		return
	}

	params, err := validateChirp(&ReqBody)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	chirpDetails, err := cfg.db.AddChirpsToUser(r.Context(), database.AddChirpsToUserParams{
		Body: params.Body,
		UserID: userID,
	})
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Something went wrong with updating chirp to user", err)
		return
	}
	respBody := NewChirpResponse(chirpDetails)
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Something went wrong with updating chirp to user", err)
		return
	} 
	
	var chirps []ChirpResponse
	for _, obj := range data {
		chirps = append(chirps, NewChirpResponse(obj))
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) 
	})
	
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerSingleChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "invalid uuid", err )
		return
	}

	data, err := cfg.db.GetSingleChirp(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	respBody := NewChirpResponse(data)
	respondWithJSON(w, http.StatusOK, respBody)

}