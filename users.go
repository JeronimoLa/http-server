package main

import (
	"encoding/json"
	"net/http"

	"github.com/jeronimoLa/http-server/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body) // reads directly from that stream instead of first buffering everything.
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Something went wrong when creating user", err)
		return
	}
	respBody := NewUserResponse(user)
	
	respondWithJSON(w, http.StatusCreated, respBody)

}

func NewUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID: 		u.ID,
		CreatedAt: 	u.CreatedAt,
		UpdatedAt: 	u.UpdatedAt,
		Email: 		u.Email,	
	}
}

func (cfg *apiConfig) handlerDeleteUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev"{
		w.WriteHeader(http.StatusForbidden)
	}
	cfg.db.DeleteAllUsers(r.Context())
}