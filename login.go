package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jeronimoLa/http-server/internal/auth"
	"github.com/jeronimoLa/http-server/internal/database"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
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
