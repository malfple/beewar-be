package auth

import (
	"gitlab.com/otqee/otqee-be/internal/auth"
	"net/http"
)

// HandleLogout handles logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// take refresh token from cookie
	refreshToken := ""
	if refreshTokenCookie, err := r.Cookie("otqee-rtoken"); err == nil {
		refreshToken = refreshTokenCookie.Value
	}

	username := auth.ValidateRefreshToken(refreshToken)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	auth.RemoveRefreshToken(refreshToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "otqee-rtoken",
		Value:    "",
		MaxAge:   0, // immediately expire the invalid refresh token
		Path:     "/api/auth",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}
