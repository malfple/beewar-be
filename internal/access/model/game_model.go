package model

// Game is a db model
type Game struct {
	ID           uint64 `json:"id"`
	Type         uint8  `json:"type"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	PlayerCount  uint8  `json:"player_count"`
	TerrainInfo  []byte `json:"terrain_info"`
	UnitInfo     []byte `json:"unit_info"`
	MapID        uint64 `json:"map_id"`
	TurnCount    int32  `json:"turn_count"`
	TurnPlayer   int8   `json:"turn_player"`
	TimeCreated  int64  `json:"time_created"`
	TimeModified int64  `json:"time_modified"`
}

// GameUser is a db model
type GameUser struct {
	ID          uint64 `json:"id"`
	GameID      uint64 `json:"game_id"`
	UserID      uint64 `json:"user_id"`
	PlayerOrder uint8  `json:"player_order"`
	FinalRank   uint8  `json:"final_rank"`
	FinalTurns  int32  `json:"final_turns"`
}
