package auth

import (
	"github.com/google/uuid"
	"time"
)

// this struct contains the username bound to the token and its expiry
type refreshTokenInfo struct {
	Username string
	ExpireAt int64
}

// 7 days expiry time
const refreshTokenExpiry = 168 * time.Hour

// maps refresh token to username
var refreshTokenStore = make(map[string]refreshTokenInfo)

// GenerateRefreshToken generates a refresh token using uuid (16-long byte array) and binds it to username
func GenerateRefreshToken(username string) string {
	token := uuid.New().String()
	refreshTokenStore[token] = refreshTokenInfo{
		Username: username,
		ExpireAt: time.Now().Add(refreshTokenExpiry).Unix(),
	}
	return token
}

// RemoveRefreshToken removes the refresh token regardless of whether it's expired
func RemoveRefreshToken(refreshToken string) {
	delete(refreshTokenStore, refreshToken)
}

// GetUsernameFromRefreshToken returns username or empty string if token not found / expired
func GetUsernameFromRefreshToken(refreshToken string) string {
	if tokenInfo, ok := refreshTokenStore[refreshToken]; ok {
		if tokenInfo.ExpireAt <= time.Now().Unix() { // token expired
			delete(refreshTokenStore, refreshToken)
			return ""
		}
		return tokenInfo.Username
	}
	return ""
}
