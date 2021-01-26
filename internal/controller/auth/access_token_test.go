package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWT(t *testing.T) {
	userID := uint64(69)
	username := "some_username"

	token := GenerateJWT(userID, username)

	claimedUserID, claimedUsername, err := ValidateJWT(token)

	assert.Equal(t, nil, err)
	assert.Equal(t, "some_username", claimedUsername)
	assert.Equal(t, uint64(69), claimedUserID)
}
