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
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"` // if this field has its zero value, don’t include it in the JSON output.
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func NewUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed.Bool,
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
		Email:          params.Email,
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

func (cfg *apiConfig) handlerUpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	AccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		JSONErrorResponse(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	userID, err := auth.ValidateJWT(AccessToken, cfg.tokenSecret)
	if err != nil {
		JSONErrorResponse(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	Resp := UserParameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&Resp)
	if err != nil {
		log.Println(err)
	}

	hash, err := auth.HashPassword(Resp.Password)
	if err != nil {
		log.Println(err)
	}

	updatedUser, err := cfg.db.UpdateEmailAndPassword(r.Context(), database.UpdateEmailAndPasswordParams{
		Email:          Resp.Email,
		HashedPassword: hash,
		ID:             userID,
	})

	respondWithJSON(w, http.StatusOK, NewUserResponse(updatedUser))
}
