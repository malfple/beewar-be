package access

import (
	"database/sql"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// CreateEmptyMap creates an empty map with the specified type and size, and returns the id
func CreateEmptyMap(mapType uint8, height, width int, name string, authorUserID uint64) (uint64, error) {
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

/*
UpdateMap updates a map

only updates user-updatable fields:
 - type
 - height
 - width
 - name
 - player_count
 - terrain_info
 - unit_info
*/
func UpdateMap(id uint64, mapType uint8, height, width int, name string, playerCount uint8, terrainInfo, unitInfo []byte) error {
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
func QueryMapByID(mapID uint64) (*model.Map, error) {
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
			return nil, err
		}
		return nil, nil
	}
	return mapp, nil
}

// all map batch queries/searches are returned in descending order of id
// currently, we ues OFFSET :'), because it's easy and the speed doesn't matter because we won't have millions of maps anyway.
// feel free to speed this up

// QueryMaps gets a list of maps
func QueryMaps(limit, offset int) ([]*model.Map, error) {
	rows, err := db.Query(`SELECT * FROM map_tab ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "map_tab"), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	maps := make([]*model.Map, 0)
	for rows.Next() {
		mapp := &model.Map{}
		err := rows.Scan(
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
			logger.GetLogger().Error("db: query error", zap.String("table", "map_tab"), zap.Error(err))
			return nil, err
		}
		maps = append(maps, mapp)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().Error("db: query error", zap.String("table", "map_tab"), zap.Error(err))
		return nil, err
	}
	return maps, nil
}
