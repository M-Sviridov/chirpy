package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	getBearerTokenTests := []struct {
		name     string
		headers  http.Header
		hasToken string
		hasErr   bool
	}{
		{
			name: "Valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer tokenx10202"},
			},
			hasToken: "tokenx10202",
			hasErr:   false,
		},
		{
			name:     "Missing authorization header",
			headers:  http.Header{},
			hasToken: "",
			hasErr:   true,
		},
		{
			name: "Malformed authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer tokenx10202"},
			},
			hasToken: "",
			hasErr:   true,
		},
	}

	for _, tt := range getBearerTokenTests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)

			if (err != nil) != tt.hasErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.hasErr)
				return
			}

			if gotToken != tt.hasToken {
				t.Errorf("GetBearerToken() = %q, want %q", gotToken, tt.hasToken)
			}
		})
	}
}
