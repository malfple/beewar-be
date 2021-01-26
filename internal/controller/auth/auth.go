package auth

import (
	"gitlab.com/beewar/beewar-be/internal/access"
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

// Register registers a new user with the provided credentials and returns status as err
func Register(email, username, password string) error {
	// for now, this is a somewhat useless function
	// need to add validations for duplicate email and username and prevent access layer from handling it
	return access.CreateUser(email, username, password)
}
