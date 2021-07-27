package mapmanager

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/formatter"
)

var (
	errDB            = errors.New("something is wrong with the database")
	errMapNotFound   = errors.New("map not found")
	errNotAuthor     = errors.New("only author can update map")
	errMapWidth      = errors.New("width must be at least 1 and at most 50")
	errMapHeight     = errors.New("height must be at least 1 and at most 50")
	errMapNameLength = errors.New("name must be at most 255")
)

const (
	mapMaxHeight     = 50
	mapMaxWidth      = 50
	mapMaxNameLength = 255
)

// CreateEmptyMap creates an empty map of fixed size and name
func CreateEmptyMap(userID uint64) (uint64, error) {
	mapID, err := access.CreateEmptyMap(0, 10, 10, "Untitled", userID, false)
	if err != nil {
		return 0, err
	}
	return mapID, nil
}

// UpdateMap updates the map
func UpdateMap(userID uint64, mapID uint64, mapType uint8, height, width int, name string, playerCount uint8, terrainInfo, unitInfo []byte) error {
	mapModel, err := access.QueryMapByID(mapID)
	if err != nil {
		return errDB
	}
	if mapModel == nil {
		return errMapNotFound
	}
	if userID != mapModel.AuthorUserID {
		return errNotAuthor
	}

	if height < 1 || height > mapMaxHeight {
		return errMapHeight
	}
	if width < 1 || width > mapMaxWidth {
		return errMapWidth
	}
	if len(name) > mapMaxNameLength {
		return errMapNameLength
	}
	if err := formatter.ValidateTerrainInfo(height, width, terrainInfo); err != nil {
		return err
	}
	if err := formatter.ValidateUnitInfo(height, width, int(playerCount), unitInfo, false); err != nil {
		return err
	}

	return access.UpdateMap(mapID, mapType, height, width, name, playerCount, terrainInfo, unitInfo)
}
