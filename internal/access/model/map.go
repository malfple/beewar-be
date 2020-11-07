package model

// Map is a db model
type Map struct {
	ID            int64  `json:"id"`
	Type          int8   `json:"type"`
	Width         int8   `json:"width"`
	Height        int8   `json:"height"`
	Name          string `json:"name"`
	PlayerCount   int8   `json:"player_count"`
	TerrainInfo   []byte `json:"terrain_info"`
	UnitInfo      []byte `json:"unit_info"`
	AuthorUserID  int64  `json:"author_user_id"`
	StatVotes     int32  `json:"stat_votes"`
	StatPlayCount int32  `json:"stat_play_count"`
	TimeCreated   int64  `json:"time_created"`
	TimeModified  int64  `json:"time_modified"`
}
