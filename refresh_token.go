package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jeronimoLa/http-server/internal/auth"
	"github.com/jeronimoLa/http-server/internal/database"
)

type RefreshTokenReponse struct {
	AccessToken string `json:"token"`
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
