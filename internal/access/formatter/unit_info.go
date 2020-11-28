package formatter

import (
	"errors"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects/units"
)

/*
terrain info and map info format
let W = map width, H = map height

--- map info ---

let N = number of units
size: variable, depends on the units
format:

N is stored implicitly -> the units are read until end of array,
and has to match their sum of byte requirements without excess information

each unit have their own requirements depending on its type, but most need 5 bytes
y, x, p, t, f
y = row number, x = column number
p = the player who owns this unit (0..number of players)
t = unit type
f = flags (unit state)

no two units can share the same position

the most basic information that a flag has are:
- life point - most units have 10, so that's only 4 bits
- turn state - most units only have two states: not yet moved, and already moved: hence 1 bit
	- this complicates more when the units store firing states and can carry other units

*/

var unitCapReq = make(map[int8]int8)

func init() {
	// --- per unit requirements ---

	// You is typical
	unitCapReq[units.UnitTypeYou] = 5

	// Infantry is typical
	unitCapReq[units.UnitTypeInfantry] = 5
}

// ErrMapInvalidUnitInfo is returned when unit info does not follow format
var ErrMapInvalidUnitInfo = errors.New("invalid unit info")

// ValidateUnitInfo validates whether unit info follows format
func ValidateUnitInfo(width, height uint8, unitInfo []byte) error {
	if len(unitInfo)%5 != 0 {
		return ErrMapInvalidUnitInfo
	}
	return nil
}

// ModelToGameUnit converts unit info from model.Game to objects.Game
func ModelToGameUnit(width, height uint8, unitInfo []byte) [][]*units.Unit {
	_units := make([][]*units.Unit, height)
	for i := uint8(0); i < height; i++ {
		_units[i] = make([]*units.Unit, width)
		for j := uint8(0); j < width; j++ {
			_units[i][j] = nil
		}
	}
	// TODO: translate unit and assign to cells
	return _units
}

// GameUnitToModel converts unit info from objects.Game to model.Game
func GameUnitToModel(width, height int8, _units [][]*units.Unit) []byte {
	return nil
}