package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var SecretKey = []byte("supersecret")

func GetSHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

type CustomClaims struct {
	Username  string   `json:"username"`
	Endpoints []string `json:"endpoints"`
	jwt.RegisteredClaims
}

func CreateJWTToken(username string, endpoints []string, secret []byte) (string, error) {
	claims := CustomClaims{
		Username:  username,
		Endpoints: endpoints,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
