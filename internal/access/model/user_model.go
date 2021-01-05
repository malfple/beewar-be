package model

// User is a db model
type User struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	Rating      uint16 `json:"rating"`
	MovesMade   uint64 `json:"moves_made"`
	GamesPlayed uint32 `json:"games_played"`
	TimeCreated int64  `json:"time_created"`
}
