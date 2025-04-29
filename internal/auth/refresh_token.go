package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error){

	bytes := make([]byte, 32)
	rand.Read(bytes)
	string := hex.EncodeToString(bytes)
	return string, nil
}
