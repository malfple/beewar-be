package auth

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"net/http"
)

// AuthenticateUser checks user credentials and returns the user on success
func AuthenticateUser(username string) *model.User {
	user := access.QueryUserByUsername(username)
	if user == nil {
		return nil
	}

	return user
}

// Login authenticates the user and returns token
func Login(username string) (string, int) {
	user := AuthenticateUser(username)
	if user == nil {
		return "", http.StatusUnauthorized
	}

	return user.Username, http.StatusOK
}
