package formatter

import "errors"

/*
terrain info and map info format
let W = map width, H = map height

--- terrain info ---

size = W * H
ordered by row first (H) then column (W)


--- arrangement ---

the cells are arranged in hexagonal cells (even-r horizontal layout)

a coordinate (y, x) indicates a cell at row y and column x.
rows and columns start at 0

even-r horizontal layout ---
this means all cells in the same row are connected straightly, but all cells in the same column are staggered,
with the even rows shifted a bit to the right

this means for a cell (y, x) in an even row, it is adjacent to
(y, x-1) (y, x+1) (y-1, x) (y+1, x)   (y-1, x+1) (y+1, x+1)
and for a cell in odd rows,
(y, x-1) (y, x+1) (y-1, x) (y+1, x)   (y-1, x-1) (y+1, x-1)

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
