package main

import (
	"sync/atomic"
	"github.com/jeronimoLa/http-server/internal/database"
)

type apiConfig struct {
	platform       string
	db             *database.Queries
	fileserverHits atomic.Int32
}
