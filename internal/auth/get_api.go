package auth

import (
	"fmt"
	"net/http"
	"strings"
)



func GetApiKey(headers http.Header) (string, error) {

	rawAuthHeader := headers.Get("Authorization")
	if rawAuthHeader == "" {
		return "", fmt.Errorf("Authorization field in header is empty")
	}
	headerVals := strings.Split(rawAuthHeader, " ")
	if len(headerVals) != 2 {
		return "", fmt.Errorf("Authorization field does not have the expected amount of values")
	}
	if headerVals[0] != "ApiKey" {

		return "", fmt.Errorf("Authorization type is not ApiKey")
	}
	return headerVals[1], nil
}
