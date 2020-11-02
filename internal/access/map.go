package access

import (
	"database/sql"
	"errors"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

var (
	// ErrMapWidth is returned when width is out of constraint
	ErrMapWidth = errors.New("width must be at least 1 and at most 50")
	// ErrMapHeight is returned when height is out of constraint
	ErrMapHeight = errors.New("height must be at least 1 and at most 50")
	// ErrMapNameLength is returned when name exceeds maximum length
	ErrMapNameLength = errors.New("name must be at most 255")
)

const (
	mapMaxWidth      = 50
	mapMaxHeight     = 50
	mapMaxNameLength = 255
)

// CreateEmptyMap creates an empty map with the specified type and size, and returns the id
func CreateEmptyMap(mapType, width, height int8, name string, authorUserID int64) (int64, error) {
	if width < 1 || width > mapMaxWidth {
		return 0, ErrMapWidth
	}
	if height < 1 || height > mapMaxHeight {
		return 0, ErrMapHeight
	}
	if len(name) > mapMaxNameLength {
		return 0, ErrMapNameLength
	}

	terrainInfo := make([]byte, width*height)
	unitInfo := make([]byte, 1)

	res, err := db.Exec(`INSERT INTO map_tab(type, width, height, name, terrain_info, unit_info, author_user_id, time_created, time_modified) VALUES (?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`,
		mapType, width, height, name, terrainInfo, unitInfo, authorUserID)
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
		&mapp.Name,
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
