package main

import (
	"database/sql"
	"encoding/json"
	"errors"

	// "fmt"
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

type RefreshTokenReponse struct {
	AccessToken string `json:"token"`
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
	var expiresIn time.Duration
	defaultExpirationTime := 3600 * time.Second
	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > int(defaultExpirationTime.Seconds()) {
		expiresIn = defaultExpirationTime
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(data.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "token was not created, try again", err)
	}

	refresh_token, err := auth.MakeRefreskToken()
	if err != nil {
		log.Println(err)
	}

	refreshTokenExists, err := cfg.db.GetRefreshTokenByEmail(r.Context(), data.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
				Token:     refresh_token,
				UpdatedAt: time.Now(),
				Email:     data.Email,
				UserID:    data.ID,
				ExpiresAt: time.Now().AddDate(0, 0, 60),
			})
		} else {
			log.Println(err)
		}

	} else {
		cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
			UpdatedAt: time.Now(),
			RevokedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Token: refreshTokenExists.Token,
		})
		cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
			Token:     refresh_token,
			UpdatedAt: time.Now(),
			Email:     data.Email,
			UserID:    data.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60),
		})
	}

	respBody := NewUserResponse(data)
	respBody.Token = token
	respBody.RefreshToken = refresh_token
	respondWithJSON(w, http.StatusOK, respBody)

}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "refresh token is not valid", err)
		return
	}

	if user.RevokedAt.Valid {
		JSONErrorResponse(w, http.StatusUnauthorized, "refresh token has been revoked ", err)
		return
	}

	expiresIn := 3600 * time.Second

	NewAccessToken, err := auth.MakeJWT(user.UserID, cfg.tokenSecret, expiresIn)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "token was not created, try again", err)
	}

	respondWithJSON(w, http.StatusOK, RefreshTokenReponse{
		AccessToken: NewAccessToken,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: refreshToken,
	})
	w.WriteHeader(204)
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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	ChirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println(err)
	}

	AccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		JSONErrorResponse(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	UserID, err := auth.ValidateJWT(AccessToken, cfg.tokenSecret)
	if err != nil {
		JSONErrorResponse(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), ChirpID)
	if errors.Is(err, sql.ErrNoRows) {
		JSONErrorResponse(w, http.StatusNotFound, err.Error(), err)
	}

	if UserID != chirp.UserID {
		JSONErrorResponse(w, http.StatusForbidden, "", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, cfg.db.DeleteSingleChirp(r.Context(), ChirpID))

}
