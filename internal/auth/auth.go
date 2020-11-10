package auth

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Login authenticates the user and returns token
func Login(username, password string) (string, int) {
	user := access.QueryUserByUsername(username)
	if user == nil {
		return "", http.StatusUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", http.StatusUnauthorized
	}

	return user.Username, http.StatusOK
}
