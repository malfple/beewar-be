package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"time"
)

// we use jwt for access token

var jwtSecret = []byte{184, 208, 147, 205, 37, 218, 186, 230, 51, 67, 100, 192, 190, 207, 108, 26, 136, 235, 173, 57, 198, 159, 15, 75, 166, 148, 108, 239, 12, 77, 164, 9}

const jwtExpiry = 15 * time.Minute

// JWTCustomClaim is a custom jwt claim
// username will be stored in `sub`
type JWTCustomClaim struct {
	UserID uint64 `json:"user_id,omitempty"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT using userID and username as claim
func GenerateJWT(userID uint64, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaim{
		userID,
		jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: time.Now().Add(jwtExpiry).Unix(),
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.GetLogger().Error("jwt signing error", zap.Error(err))
	}

	return tokenString
}

// ValidateJWT returns the userID and username. If it's not valid, an error is returned
func ValidateJWT(tokenString string) (uint64, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(*JWTCustomClaim); ok && token.Valid {
		return claims.UserID, claims.Subject, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}
