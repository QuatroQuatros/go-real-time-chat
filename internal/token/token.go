package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("super-secret-key")

type Claims struct {
	UserID uint `json:"userId"`
	Guest  bool `json:"guest"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, guest bool) (string, error) {
	claims := Claims{
		UserID: userID,
		Guest:  guest,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
