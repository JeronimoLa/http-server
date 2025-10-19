package auth

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"log"
	"net/http"
	"strings"
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

	trimmedKey, ok := strings.CutPrefix(apiKey, "ApiKey")
	if !ok {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}
	cleanApiKey := strings.TrimSpace(trimmedKey)

	return cleanApiKey, nil
}
