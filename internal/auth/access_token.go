package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"time"
)

// we use jwt for access token

var jwtSecret = []byte{184, 208, 147, 205, 37, 218, 186, 230, 51, 67, 100, 192, 190, 207, 108, 26, 136, 235, 173, 57, 198, 159, 15, 75, 166, 148, 108, 239, 12, 77, 164, 9}

const jwtExpiry = 15 * time.Minute

// GenerateJWT generates a JWT using username as claim
func GenerateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(jwtExpiry).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.GetLogger().Error("jwt signing error", zap.Error(err))
	}

	return tokenString
}

// ValidateJWT returns the username. If it's not valid, an error is returned
func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", fmt.Errorf("invalid token")
}
