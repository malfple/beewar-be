package model

// User is a db model
type User struct {
	ID          int64
	Email       string
	Username    string
	Password    string
	TimeCreated int64
}
