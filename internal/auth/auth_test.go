package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidJWT(t *testing.T) {
	user := uuid.New()
	secret := "secret-a"

	jwt, err := MakeJWT(user, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	validUser, err := ValidateJWT(jwt, secret)
	if err != nil {
		t.Fatal(err)
	}

	if validUser != user {
		t.Error("The returned User Does Not Matche The Current User")
	}
}

func TestExpiredToken(t *testing.T) {
	user := uuid.New()
	secret := "secret-b"

	jwt, err := MakeJWT(user, secret, time.Duration(-1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwt, secret)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestWrongSecret(t *testing.T) {
	user := uuid.New()
	secret := "secret-c"

	jwt, err := MakeJWT(user, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwt, "secret-d")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
