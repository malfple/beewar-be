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

func TestMultipleRefreshTokens(t *testing.T) {
	userID := uint64(465)
	username := "some_other_username"
	token1 := GenerateRefreshToken(userID, username)
	token2 := GenerateRefreshToken(userID, username)
	token3 := GenerateRefreshToken(userID, username)

	userID1, _ := ValidateRefreshToken(token1)
	assert.Equal(t, uint64(0), userID1)
	userID2, _ := ValidateRefreshToken(token2)
	assert.Equal(t, uint64(0), userID2)
	userID3, _ := ValidateRefreshToken(token3)
	assert.Equal(t, userID, userID3)
}
