package auth

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	// ErrUsernameLength is returned when username length is out of constraint
	ErrUsernameLength = errors.New("username length must be at least 1 and at most 255")
	// ErrPasswordLength is returned when password length is out of constraint
	ErrPasswordLength = errors.New("password length must be at least 8 and at most 255")
)

const (
	userMaxUsernameLength = 255
	userMinPasswordLength = 8
	userMaxPasswordLength = 255
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
	if len(username) < 1 || len(username) > userMaxUsernameLength {
		return ErrUsernameLength
	}
	if len(password) < userMinPasswordLength || len(password) > userMaxPasswordLength {
		return ErrPasswordLength
	}
	// need to add validations for duplicate email and username and prevent access layer from handling it
	return access.CreateUser(email, username, password)
}
