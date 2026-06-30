package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefreshToken(email, role string, id int64) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	return refreshToken.SignedString([]byte(secretKey))
}
