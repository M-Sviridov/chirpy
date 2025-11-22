package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "password", time.Hour)

	validateJWTTests := []struct {
		name        string
		tokenString string
		tokenSecret string
		hasUserID   uuid.UUID
		hasErr      bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "password",
			hasUserID:   userID,
			hasErr:      false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalidtokenstring",
			tokenSecret: "password",
			hasUserID:   uuid.Nil,
			hasErr:      true,
		},
		{
			name:        "Invalid token secret",
			tokenString: validToken,
			tokenSecret: "invalidtokensecret",
			hasUserID:   uuid.Nil,
			hasErr:      true,
		},
	}

	for _, tt := range validateJWTTests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.hasErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.hasErr)
			}
			if gotUserID != tt.hasUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.hasUserID)
			}
		})
	}
}
