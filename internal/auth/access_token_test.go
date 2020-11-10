package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWT(t *testing.T) {
	username := "some_username"

	token := GenerateJWT(username)

	claimedUsername, err := ValidateJWT(token)

	assert.Equal(t, nil, err)
	assert.Equal(t, "some_username", claimedUsername)
}
