package auth_test

import (
	"github.com/google/uuid"
	"github.com/ywallis/chirpy/internal/auth"
	"testing"
	"time"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		// TODO: Add test cases.
		{name: "standard",
			userID:      uuid.New(),
			tokenSecret: "password",
			expiresIn:   time.Duration(time.Second * 5),
			wantErr:     false},

		{name: "no password",
			userID:      uuid.New(),
			tokenSecret: "",
			expiresIn:   time.Duration(time.Second * 5),
			wantErr:     true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		// TODO: Add test cases.
		{name: "standard",
			userID:      uuid.New(),
			tokenSecret: "password",
			expiresIn:   time.Duration(time.Second * 5),
			wantErr:     false},

		{name: "no password",
			userID:      uuid.New(),
			tokenSecret: "",
			expiresIn:   time.Duration(time.Second * 5),
			wantErr:     true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", err)
				}
				return
			}
			returnedUuid, err := auth.ValidateJWT(token, tt.tokenSecret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
			if returnedUuid != tt.userID {
				t.Fatal("Tokens could not be reconciliated")
			}
		})
	}
}
