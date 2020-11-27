package model

// Map is a db model
type Map struct {
	ID            uint64 `json:"id"`
	Type          uint8  `json:"type"`
	Width         uint8  `json:"width"`
	Height        uint8  `json:"height"`
	Name          string `json:"name"`
	PlayerCount   uint8  `json:"player_count"`
	TerrainInfo   []byte `json:"terrain_info"`
	UnitInfo      []byte `json:"unit_info"`
	AuthorUserID  uint64 `json:"author_user_id"`
	StatPlayCount uint32 `json:"stat_play_count"`
	TimeCreated   int64  `json:"time_created"`
	TimeModified  int64  `json:"time_modified"`
}
