package message

// JoinMessageData is a message struct for client to send to server
type JoinMessageData struct {
	PlayerOrder uint8 `json:"player_order"`
}

// JoinMessageDataExt is an event message by server to notify if a player joined the game.
// This uses a struct from game data message
type JoinMessageDataExt struct {
	Player *Player `json:"player"`
}
