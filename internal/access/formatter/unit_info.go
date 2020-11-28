package formatter

import "errors"

/*
terrain info and map info format
let W = map width, H = map height

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

// ErrMapInvalidUnitInfo is returned when unit info does not follow format
var ErrMapInvalidUnitInfo = errors.New("invalid unit info")

// ValidateUnitInfo validates whether unit info follows format
func ValidateUnitInfo(width, height uint8, unitInfo []byte) error {
	if len(unitInfo)%5 != 0 {
		return ErrMapInvalidUnitInfo
	}
	return nil
}
