package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jeronimoLa/http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Println(err)
	}
	dbQueries := database.New(db)



	const port = "8080"
	const filepathRoot = "."

	apiCfg := &apiConfig{databaseQueries: dbQueries}
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := &http.Server{
		Addr:		":" + port,
		Handler: 	mux,
	}

	server.ListenAndServe()
	log.Print("Listening...")
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}	

