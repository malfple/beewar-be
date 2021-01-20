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
