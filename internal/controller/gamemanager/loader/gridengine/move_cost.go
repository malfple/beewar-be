package gridengine

import "gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"

// CalcMoveCost calculates move cost for a terrain given several factors. This only affects ground move type.
func CalcMoveCost(terrainType int, unitWeight int) int {
	switch terrainType {
	case objects.TerrainTypeVoid:
		return 999999
	case objects.TerrainTypePlains:
		return 1
	case objects.TerrainTypeWalls:
		return 999999
	case objects.TerrainTypeHoneyField:
		return 2
	case objects.TerrainTypeWasteland:
		return 1 + unitWeight
	case objects.TerrainTypeIceField:
		return 2
	case objects.TerrainTypeThrone:
		return 1
	}
	return 999999
}
