package model

// Game is a db model
type Game struct {
	ID           int64
	Type         int8
	Width        int8
	Height       int8
	PlayerCount  int8
	TerrainInfo  []byte
	UnitInfo     []byte
	MapID        int64
	TurnCount    int32
	TurnPlayer   int8
	TimeCreated  int64
	TimeModified int64
}

// GameUser is a db model
type GameUser struct {
	ID     int64
	GameID int64
	UserID int64
}
