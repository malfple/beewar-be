package model

// User is a db model
type User struct {
	ID          uint64
	Email       string
	Username    string
	Password    string
	Rating      uint16
	MovesMade   uint64
	GamesPlayed uint32
	TimeCreated int64
}
