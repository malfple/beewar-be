package model

// User is a db model
type User struct {
	ID              uint64 `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password,omitempty"`
	Rating          int32  `json:"rating"`
	MovesMade       uint64 `json:"moves_made"`
	GamesPlayed     uint32 `json:"games_played"`
	HighestCampaign int32  `json:"highest_campaign"`
	TimeCreated     int64  `json:"time_created"`
}
