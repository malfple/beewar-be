package access

import (
	"database/sql"
	"errors"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
)

/*
terrain info and map info format
let W = map width, H = map height

--- terrain info ---

size = W * H
ordered by row first (H) then column (W)

--- map info ---

let N = number of units
size: 5 * N
format:
for each unit, there are 5 numbers
y, x, p, t, f
y = row number, x = column number
p = the player who owns this unit (0..number of players)
t = unit type
f = flags (unit state)

no two units can share the same position

*/

var (
	// ErrMapWidth is returned when width is out of constraint
	ErrMapWidth = errors.New("width must be at least 1 and at most 50")
	// ErrMapHeight is returned when height is out of constraint
	ErrMapHeight = errors.New("height must be at least 1 and at most 50")
	// ErrMapNameLength is returned when name exceeds maximum length
	ErrMapNameLength = errors.New("name must be at most 255")
	// ErrMapInvalidTerrainInfo is returned when terrain info does not match map size
	ErrMapInvalidTerrainInfo = errors.New("invalid terrain info")
	// ErrMapInvalidUnitInfo is returned when unit info does not follow format
	ErrMapInvalidUnitInfo = errors.New("invalid unit info")
)

const (
	mapMaxWidth      = 50
	mapMaxHeight     = 50
	mapMaxNameLength = 255
)

// TODO: complete validations
func validateTerrainInfo(width, height int8, terrainInfo []byte) error {
	if len(terrainInfo) != int(width)*int(height) {
		return ErrMapInvalidTerrainInfo
	}
	return nil
}

func validateUnitInfo(width, height int8, unitInfo []byte) error {
	if len(unitInfo)%5 != 0 {
		return ErrMapInvalidUnitInfo
	}
	return nil
}

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
	unitInfo := make([]byte, 0)

	const stmtCreateEmptyMap = `INSERT INTO map_tab
(type, width, height, name, player_count, terrain_info, unit_info, author_user_id, time_created, time_modified)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())`

	res, err := db.Exec(stmtCreateEmptyMap,
		mapType, width, height, name, 1, terrainInfo, unitInfo, authorUserID)
	if err != nil {
		logger.GetLogger().Error("db: insert error", zap.String("table", "map_tab"), zap.Error(err))
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateMap updates a map
func UpdateMap(id int64, mapType, width, height int8, name string, playerCount int8, terrainInfo, unitInfo []byte) error {
	if width < 1 || width > mapMaxWidth {
		return ErrMapWidth
	}
	if height < 1 || height > mapMaxHeight {
		return ErrMapHeight
	}
	if len(name) > mapMaxNameLength {
		return ErrMapNameLength
	}
	if err := validateTerrainInfo(width, height, terrainInfo); err != nil {
		return err
	}
	if err := validateUnitInfo(width, height, unitInfo); err != nil {
		return err
	}

	const stmtUpdateMap = `UPDATE map_tab
SET type=?, width=?, height=?, name=?, player_count=?, terrain_info=?, unit_info=?
WHERE id=?`

	_, err := db.Exec(stmtUpdateMap,
		mapType, width, height, name, playerCount, terrainInfo, unitInfo,
		id)
	if err != nil {
		logger.GetLogger().Error("db: update error", zap.String("table", "map_tab"), zap.Error(err))
		return err
	}

	return nil
}

// QueryMapByID gets a single map by id
func QueryMapByID(mapID int64) *model.Map {
	row := db.QueryRow(`SELECT * FROM map_tab WHERE id=? LIMIT 1`, mapID)

	mapp := &model.Map{}
	err := row.Scan(
		&mapp.ID,
		&mapp.Type,
		&mapp.Width,
		&mapp.Height,
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
