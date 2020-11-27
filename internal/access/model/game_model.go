package model

// Game is a db model
type Game struct {
	ID           uint64
	Type         uint8
	Width        uint8
	Height       uint8
	PlayerCount  uint8
	TerrainInfo  []byte
	UnitInfo     []byte
	MapID        uint64
	TurnCount    int32
	TurnPlayer   int8
	TimeCreated  int64
	TimeModified int64
}

// GameUser is a db model
type GameUser struct {
	ID          uint64
	GameID      uint64
	UserID      uint64
	RankOrder   uint8
	TurnsLasted int32
}
