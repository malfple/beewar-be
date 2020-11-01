package model

// User is a db model
type User struct {
	ID          int64
	Email       string
	Username    string
	Password    string
	Rating      int16
	MovesMade   int64
	GamesPlayed int32
	TimeCreated int64
}
