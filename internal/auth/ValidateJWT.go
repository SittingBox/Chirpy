package auth

import (
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil(), err
	}

	tokenSubject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	uniqueUserID, err := uuid.Parse(tokenSubject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uniqueUserID, nil

}
