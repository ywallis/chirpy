package auth_test

import (
	"github.com/ywallis/chirpy/internal/auth"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		wantErr  bool
	}{

		{name: "6char", password: "abcdef", wantErr:  false},
		{name: "4char", password: "abcd", wantErr:  false},
		{name: "empty password", password: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}
			if hash == "" {
				t.Error("Hashed password should not be empty")
			}
		})
	}
}
