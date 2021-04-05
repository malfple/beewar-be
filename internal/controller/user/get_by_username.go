package user

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
)

// GetByUsername returns user with password hidden.
func GetByUsername(username string) *model.User {
	user := access.QueryUserByUsername(username)
	if user != nil {
		user.Password = ""
	}
	return user
}
