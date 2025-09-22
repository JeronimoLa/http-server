package main

import (
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("GET /admin/metrics", apiCfg.ReadMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	server := &http.Server{
		Addr:		":" + port,
		Handler: 	mux,
	}

	server.ListenAndServe()
	log.Print("Listening...")
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}	

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func ValidateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters %s", err)
		w.WriteHeader(500)
		return
	}

	
	if len(params.Body) > 140{	
		type returnVals struct {
			ErrorMessage 	string `json:"error"`
		}

		respBody := returnVals{
			ErrorMessage: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
	} else {

	
		type SuccessReturnVals struct {
				SuccessMessage 	bool `json:"valid"`
		}
		respBody := SuccessReturnVals{
			SuccessMessage: true,
		}
		data, err := json.Marshal(respBody)

		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
	}



}


func (cfg *apiConfig) ReadMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	
	// w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
	w.Write([]byte(fmt.Sprintf(`<html>
									<body>
										<h1>Welcome, Chirpy Admin</h1>
										<p>Chirpy has been visited %d times!</p>
									</body>
								</html>`, cfg.fileserverHits.Load())))

	// fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}