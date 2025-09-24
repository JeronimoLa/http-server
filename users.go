package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jeronimoLa/http-server/internal/auth"
	"github.com/jeronimoLa/http-server/internal/database"

)

type UserParameters struct {
		Password	string `json:"password"`
		Email 		string `json:"email"`
	}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func NewUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body) // reads directly from that stream instead of first buffering everything.
	params := UserParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println(err)
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Something went wrong when creating user", err)
		return
	}
	respBody := NewUserResponse(user)
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (cfg *apiConfig) handlerDeleteUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
	}
	cfg.db.DeleteAllUsers(r.Context())
}

func (cfg *apiConfig) HanderLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := UserParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	data, err := cfg.db.GetPasswordByEmail(r.Context(), params.Email)
	if err != nil {
		JSONErrorResponse(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, data.HashedPassword)
	if err != nil {
	    JSONErrorResponse(w, http.StatusInternalServerError, "server error", err)
		return
	} 
	if !match {
		JSONErrorResponse(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	respBody := NewUserResponse(data)
	respondWithJSON(w, http.StatusOK, respBody)

}
