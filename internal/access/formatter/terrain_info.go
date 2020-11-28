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
