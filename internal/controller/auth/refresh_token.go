package auth

import (
	"github.com/google/uuid"
	"time"
)

// we use uuid for refresh token

// RefreshTokenCookieName is the cookie name that stores the refresh token
const RefreshTokenCookieName = "beewar-rtoken"

// 7 days expiry time
const refreshTokenExpiry = 168 * time.Hour

// this struct contains the username bound to the token and its expiry
type refreshTokenInfo struct {
	UserID   uint64
	Username string
	ExpireAt int64
}

// maps refresh token to toke infos: refreshTokenInfo
var refreshTokenStore = make(map[string]refreshTokenInfo)

// maps userID to refresh token
var userIDToRefreshTokenMap = make(map[uint64]string)

// GenerateRefreshToken generates a refresh token using uuid (16-long byte array) and binds it to username/userid
func GenerateRefreshToken(userID uint64, username string) string {
	token := uuid.New().String()
	refreshTokenStore[token] = refreshTokenInfo{
		UserID:   userID,
		Username: username,
		ExpireAt: time.Now().Add(refreshTokenExpiry).Unix(),
	}
	if oldToken, ok := userIDToRefreshTokenMap[userID]; ok { // user already has refresh token
		RemoveRefreshToken(oldToken)
	}
	userIDToRefreshTokenMap[userID] = token
	return token
}

// RemoveRefreshToken removes the refresh token regardless of whether it's expired
func RemoveRefreshToken(refreshToken string) {
	delete(refreshTokenStore, refreshToken)
}

// ValidateRefreshToken checks refresh token and returns username,
// or empty string if token not found / expired
func ValidateRefreshToken(refreshToken string) (uint64, string) {
	if tokenInfo, ok := refreshTokenStore[refreshToken]; ok {
		if time.Now().Unix() > tokenInfo.ExpireAt { // token expired
			delete(refreshTokenStore, refreshToken)
			return 0, ""
		}
		return tokenInfo.UserID, tokenInfo.Username
	}
	return 0, ""
}

// GetRefreshTokenCount returns the number of active refresh tokens (or in other words, session)
func GetRefreshTokenCount() int {
	return len(refreshTokenStore)
}
