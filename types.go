package main

import (
	"sync/atomic"

	"github.com/jeronimoLa/http-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databaseQueries *database.Queries
}



