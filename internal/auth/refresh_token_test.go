package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	username := "some_other_username"
	token := GenerateRefreshToken(username)

	username1 := ValidateRefreshToken(token)
	assert.Equal(t, username, username1)

	RemoveRefreshToken(token)
	username2 := ValidateRefreshToken(token)
	assert.Equal(t, "", username2)
}
