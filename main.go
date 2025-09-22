package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := &apiConfig{}
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

