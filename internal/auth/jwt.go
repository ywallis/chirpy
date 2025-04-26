package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if tokenSecret == "" {
		return "", fmt.Errorf("Token secret should not be empty")
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := jwt.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	userUuid, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userUuid, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	rawAuthHeader := headers.Get("Authorization")
	if rawAuthHeader == "" {
		return "", fmt.Errorf("Authorization field in header is empty")
	}
	headerVals := strings.Split(rawAuthHeader, " ")
	if len(headerVals) != 2 {
		return "", fmt.Errorf("Authorization field does not have the expected amount of values")
	}
	if headerVals[0] != "Bearer" {

		return "", fmt.Errorf("Authorization type is not bearer")
	}
	return headerVals[1], nil
}
