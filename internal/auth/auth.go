package auth

import (
	"log"
	"net/http"
	"fmt"
	"strings"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println(err)
	}
	return hashed_password, nil

}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no token provided")	
	}

	trimmedStr, ok := strings.CutPrefix(apiKey, "ApiKey")
	if !ok {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}
	cleanedApiKey := strings.TrimSpace(trimmedStr)
	
	return cleanedApiKey, nil
}
