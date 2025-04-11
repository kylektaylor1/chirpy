package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	pw := "password"
	hashedPw, err := HashPassword(pw)
	if err != nil {
		t.Errorf("Unexpected error hashing password. Err: %v", err)
	}
	chErr := CheckPasswordHash(hashedPw, pw)
	if chErr != nil {
		t.Errorf("Password should be valid. Err: %v", chErr)
	}
}

func TestCanMakeJwt(t *testing.T) {
	userId := uuid.New()
	secret := "some-secret"
	_, err := MakeJWT(userId, secret, time.Hour)
	if err != nil {
		t.Errorf("cannot make jwt -- %v\n", err)
	}
}

func TestValidateNormalJwt(t *testing.T) {
	userId := uuid.New()
	secret := "some-secret"
	jwt, err := MakeJWT(userId, secret, time.Hour)
	if err != nil {
		t.Errorf("cannot make jwt -- %v\n", err)
	}
	_, validErr := ValidateJWT(jwt, secret)
	if validErr != nil {
		t.Errorf("should be valid -- %v\n", validErr)
	}
}

func TestGetAuthHeader(t *testing.T) {
	h := make(map[string][]string, 5)
	h["Authorization"] = []string{"Bearer token"}
	fmt.Printf("h: %+v\n", h)
	_, err := GetBearerToken(h)
	if err != nil {
		t.Errorf("error getting bearer token -- %v\n", err)
	}
}
