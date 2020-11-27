package handler

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

// HandleProfile handles profile query
func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	user := access.QueryUserByUsername(username)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &ProfileResponse{User: nil}
	if user != nil {
		resp.User = &userProfile{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			Rating:      user.Rating,
			MovesMade:   user.MovesMade,
			GamesPlayed: user.GamesPlayed,
			TimeCreated: user.TimeCreated,
		}
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// ProfileResponse is a response for profile handler
type ProfileResponse struct {
	User *userProfile `json:"user"`
}

type userProfile struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Rating      uint16 `json:"rating"`
	MovesMade   uint64 `json:"moves_made"`
	GamesPlayed uint32 `json:"games_played"`
	TimeCreated int64  `json:"time_created"`
}
