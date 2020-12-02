package formatter

import "errors"

/*
terrain info and map info format
let W = map width, H = map height

--- terrain info ---

size = W * H
ordered by row first (H) then column (W)

*/

// ErrMapInvalidTerrainInfo is returned when terrain info does not match map size
var ErrMapInvalidTerrainInfo = errors.New("invalid terrain info")

// ValidateTerrainInfo validates whether a terrain info follows format
func ValidateTerrainInfo(width, height uint8, terrainInfo []byte) error {
	if len(terrainInfo) != int(width)*int(height) {
		return ErrMapInvalidTerrainInfo
	}
	return nil
}

// ModelToGameTerrain converts terrain info from model.Game to loader.GameLoader
func ModelToGameTerrain(width, height uint8, terrainInfo []byte) [][]uint8 {
	terrain := make([][]uint8, height)
	for i := uint8(0); i < height; i++ {
		terrain[i] = make([]uint8, width)
		for j := uint8(0); j < width; j++ {
			terrain[i][j] = terrainInfo[i*width+j]
		}
	}
	return terrain
}

// GameTerrainToModel converts terrain info from loader.GameLoader to model.Game
func GameTerrainToModel(width, height uint8, terrain [][]uint8) []byte {
	terrainInfo := make([]byte, int(width)*int(height))
	for i := uint8(0); i < height; i++ {
		for j := uint8(0); j < width; j++ {
			terrainInfo[i*width+j] = terrain[i][j]
		}
	}
	return terrainInfo
}
