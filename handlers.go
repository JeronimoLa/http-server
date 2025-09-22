package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request){
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