package model

// Map is a db model
type Map struct {
	ID            int64
	Type          int8
	Width         int8
	Height        int8
	TerrainInfo   []byte
	UnitInfo      []byte
	AuthorUserID  int64
	StatVotes     int32
	StatPlayCount int32
	TimeCreated   int64
	TimeModified  int64
}
