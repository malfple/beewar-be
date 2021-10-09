package formatter

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
)

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

var (
	errMapInvalidTerrainInfo = errors.New("invalid terrain info, maybe the map size does not match the contents")
	errMapUnknownTerrainType = errors.New("unknown terrain type")
	errThroneCount           = errors.New("throne terrain count does not match")
)

// ValidateTerrainInfo validates whether a terrain info follows format
func ValidateTerrainInfo(height, width int, terrainInfo []byte, expectedThroneCount int) error {
	if len(terrainInfo) != height*width {
		return errMapInvalidTerrainInfo
	}
	throneCount := 0
	for _, terrainType := range terrainInfo {
		switch terrainType {
		case objects.TerrainTypeVoid:
		case objects.TerrainTypePlains:
		case objects.TerrainTypeWalls:
		case objects.TerrainTypeHoneyField:
		case objects.TerrainTypeWasteland:
		case objects.TerrainTypeIceField:
		case objects.TerrainTypeThrone:
			throneCount++
		default:
			return errMapUnknownTerrainType
		}
	}
	if expectedThroneCount != -1 && throneCount != expectedThroneCount {
		return errThroneCount
	}
	return nil
}

// ModelToGameTerrain converts terrain info from model.Game to loader.GameLoader
func ModelToGameTerrain(height, width int, terrainInfo []byte) [][]int {
	terrain := make([][]int, height)
	for i := 0; i < height; i++ {
		terrain[i] = make([]int, width)
		for j := 0; j < width; j++ {
			terrain[i][j] = int(terrainInfo[i*width+j])
		}
	}
	return terrain
}

// GameTerrainToModel converts terrain info from loader.GameLoader to model.Game
func GameTerrainToModel(height, width int, terrain [][]int) []byte {
	terrainInfo := make([]byte, height*width)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			terrainInfo[i*width+j] = byte(terrain[i][j])
		}
	}
	return terrainInfo
}
