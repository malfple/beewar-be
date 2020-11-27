package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	userID := uint64(465)
	username := "some_other_username"
	token := GenerateRefreshToken(userID, username)

	userID1, username1 := ValidateRefreshToken(token)
	assert.Equal(t, username, username1)
	assert.Equal(t, userID, userID1)

	RemoveRefreshToken(token)
	userID2, username2 := ValidateRefreshToken(token)
	assert.Equal(t, "", username2)
	assert.Equal(t, uint64(0), userID2)
}
