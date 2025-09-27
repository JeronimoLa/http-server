package auth

import (
	"testing"
	"time"
	

	"github.com/google/uuid"

)

func TestJWTCreation(t *testing.T) {

	tests := []struct {
		name		string
		userID 		string
		tokenSecret string
		expiresIn 	time.Duration
		wantErr		bool
		// matchUserID	bool

	}{
		{
			name:			"Valid userID",
			userID:			"f0f05cf3-57ab-4783-bd37-d15641e2023a",
			tokenSecret:	"secret-token",
			expiresIn: 		10 * time.Second,
			wantErr: 		false,

		},
		{
			name:			"Invalid userID",
			userID:			"some-uuid",
			tokenSecret:	"secret-token",
			expiresIn: 		10 * time.Second,
			wantErr: 		true,

		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userUUID, err := uuid.Parse(tt.userID)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("unexpected parse error: %v", err)
				}
				return
			}

			token, err := MakeJWT(userUUID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
    			return
			}
			if token == "" && !tt.wantErr {
			    t.Errorf("expected a token, got empty string")
			}
		})
	}
}

