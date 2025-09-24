package main

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jeronimoLa/http-server/internal/database"
)

type apiConfig struct {
	platform 	string
	db 			*database.Queries
	fileserverHits atomic.Int32
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type ChirpResponse struct {
	ChirpID   uuid.UUID `json:"chirp_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	ID        uuid.UUID `json:"user_id"`
}

