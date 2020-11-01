package access

import (
	"database/sql"
	"errors"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

var (
	// ErrWidth is returned when width is out of constraint
	ErrWidth = errors.New("width must be at least 1 and at most 50")
	// ErrHeight is returned when height is out of constraint
	ErrHeight = errors.New("height must be at least 1 and at most 50")
)

const (
	// MapMaxWidth defines the maximum map width
	MapMaxWidth = 50
	// MapMaxHeight defines the maximum map height
	MapMaxHeight = 50
)

// CreateEmptyMap creates an empty map with the specified type and size, and returns the id
func CreateEmptyMap(mapType, width, height int8, authorUserID int64) (int64, error) {
	if width < 1 || width > MapMaxWidth {
		return 0, ErrWidth
	}
	if height < 1 || height > MapMaxHeight {
		return 0, ErrHeight
	}

	terrainInfo := make([]byte, width*height)
	unitInfo := make([]byte, 1)

	res, err := db.Exec(`INSERT INTO map_tab(type, width, height, terrain_info, unit_info, author_user_id, time_created, time_modified) VALUES (?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`,
		mapType, width, height, terrainInfo, unitInfo, authorUserID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "map_tab"), zap.Error(err))
		return 0, err
	}
	return res.LastInsertId()
}

// GetMapByID gets a single map by id
func GetMapByID(mapID int64) *model.Map {
	row := db.QueryRow(`SELECT * FROM map_tab WHERE id=? LIMIT 1`, mapID)

	mapp := &model.Map{}
	err := row.Scan(
		&mapp.ID,
		&mapp.Type,
		&mapp.Width,
		&mapp.Height,
		&mapp.TerrainInfo,
		&mapp.UnitInfo,
		&mapp.AuthorUserID,
		&mapp.StatVotes,
		&mapp.StatPlayCount,
		&mapp.TimeCreated,
		&mapp.TimeModified)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.GetLogger().Error("db: query error", zap.String("table", "map_tab"), zap.Error(err))
		}
		return nil
	}
	return mapp
}
