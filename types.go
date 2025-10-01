package main

import (
	"github.com/jeronimoLa/http-server/internal/database"
	"sync/atomic"
)

type apiConfig struct {
	platform       string
	db             *database.Queries
	tokenSecret    string
	polkaKey       string
	fileserverHits atomic.Int32
}
