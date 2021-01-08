package access

import (
	"database/sql"
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access/formatter"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// terrain info and unit info description is in formatter package

var (
	// ErrMapWidth is returned when width is out of constraint
	ErrMapWidth = errors.New("width must be at least 1 and at most 50")
	// ErrMapHeight is returned when height is out of constraint
	ErrMapHeight = errors.New("height must be at least 1 and at most 50")
	// ErrMapNameLength is returned when name exceeds maximum length
	ErrMapNameLength = errors.New("name must be at most 255")
)

const (
	mapMaxHeight     = 50
	mapMaxWidth      = 50
	mapMaxNameLength = 255
)

// CreateEmptyMap creates an empty map with the specified type and size, and returns the id
func CreateEmptyMap(mapType uint8, height, width int, name string, authorUserID uint64) (uint64, error) {
	if height < 1 || height > mapMaxHeight {
		return 0, ErrMapHeight
	}
	if width < 1 || width > mapMaxWidth {
		return 0, ErrMapWidth
	}
	if len(name) > mapMaxNameLength {
		return 0, ErrMapNameLength
	}

	terrainInfo := make([]byte, width*height)
	unitInfo := make([]byte, 0)

	const stmtCreateEmptyMap = `INSERT INTO map_tab
(type, height, width, name, player_count, terrain_info, unit_info, author_user_id, time_created, time_modified)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`

	res, err := db.Exec(stmtCreateEmptyMap,
		mapType, height, width, name, 1, terrainInfo, unitInfo, authorUserID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "map_tab"), zap.Error(err))
		return 0, err
	}
	lastInsertID, err := res.LastInsertId()
	return uint64(lastInsertID), err
}

// UpdateMap updates a map
func UpdateMap(id uint64, mapType uint8, height, width int, name string, playerCount uint8, terrainInfo, unitInfo []byte) error {
	if height < 1 || height > mapMaxHeight {
		return ErrMapHeight
	}
	if width < 1 || width > mapMaxWidth {
		return ErrMapWidth
	}
	if len(name) > mapMaxNameLength {
		return ErrMapNameLength
	}
	if err := formatter.ValidateTerrainInfo(height, width, terrainInfo); err != nil {
		return err
	}
	if err := formatter.ValidateUnitInfo(height, width, unitInfo); err != nil {
		return err
	}

	const stmtUpdateMap = `UPDATE map_tab
SET type=?, height=?, width=?, name=?, player_count=?, terrain_info=?, unit_info=?, time_modified=UNIX_TIMESTAMP()
WHERE id=?`

	_, err := db.Exec(stmtUpdateMap,
		mapType, height, width, name, playerCount, terrainInfo, unitInfo,
		id)
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "map_tab"), zap.Error(err))
		return err
	}

	return nil
}

// QueryMapByID gets a single map by id
func QueryMapByID(mapID uint64) *model.Map {
	row := db.QueryRow(`SELECT * FROM map_tab WHERE id=? LIMIT 1`, mapID)

	mapp := &model.Map{}
	err := row.Scan(
		&mapp.ID,
		&mapp.Type,
		&mapp.Height,
		&mapp.Width,
		&mapp.Name,
		&mapp.PlayerCount,
		&mapp.TerrainInfo,
		&mapp.UnitInfo,
		&mapp.AuthorUserID,
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
