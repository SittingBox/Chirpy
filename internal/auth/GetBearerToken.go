package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Access Denied, No auth token")
	}

	authSlice := strings.Fields(authHeader)

	if len(authSlice) != 2 {
		return "", fmt.Errorf("Header is wrong")
	}

	if authSlice[0] != "Bearer" {
		return "", fmt.Errorf("Bearer Header Missing")
	}

	bearerToken := authSlice[1]

	return bearerToken, nil
}
