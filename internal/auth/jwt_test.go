package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ywallis/chirpy/internal/auth"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{name: "happy",
			headers: http.Header{"Authorization": {"Bearer xxx"}},
			want:    "xxx",
			wantErr: false},
		{name: "empty",
			headers: http.Header{"Authorization": {""}},
			want:    "",
			wantErr: true},
		{name: "wrong auth",
			headers: http.Header{"Authorization": {"Basic foo"}},
			want:    "",
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
