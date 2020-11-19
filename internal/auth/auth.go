package auth

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Login authenticates the user and returns (refresh token, access token, http status code)
func Login(username, password string) (string, string, int) {
	user := access.QueryUserByUsername(username)
	if user == nil {
		return "", "", http.StatusUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", http.StatusUnauthorized
	}

	// username and password is valid
	return GenerateRefreshToken(user.ID, username), GenerateJWT(user.ID, username), http.StatusOK
}
