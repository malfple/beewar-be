package formatter

import (
	"errors"
	"gitlab.com/otqee/otqee-be/internal/access/formatter/objects"
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

each unit have their own requirements depending on its type, but most need 6 bytes
y, x, p, t, hp, s
y = row number, x = column number
p = the player who owns this unit (0..number of players)
t = unit type
hp = unit health point
s = state/flags (unit state)

no two units can share the same position

the most basic information that a flag has are:
- turn state - most units only have two states: not yet moved, and already moved: hence 1 bit
	- this complicates more when the units store firing states and can carry other units


unit type information is available in `units` package
*/

// ErrMapInvalidUnitInfo is returned when unit info does not follow format
var ErrMapInvalidUnitInfo = errors.New("invalid unit info")

// ValidateUnitInfo validates whether unit info follows format
func ValidateUnitInfo(width, height uint8, unitInfo []byte) error {
	for i := 0; i < len(unitInfo); {
		if i+5 >= len(unitInfo) { // the remaining length is less than 6 (the required minimum of a normal unit)
			return ErrMapInvalidUnitInfo
		}
		y := unitInfo[i]
		x := unitInfo[i+1]
		if y < 0 || y >= height || x < 0 || x >= width {
			return ErrMapInvalidUnitInfo
		}
		t := unitInfo[i+3]
		switch t {
		case objects.UnitTypeYou:
			i += 6
		case objects.UnitTypeInfantry:
			i += 6
		default:
			return ErrMapInvalidUnitInfo
		}
	}
	return nil
}

// ModelToGameUnit converts unit info from model.Game to loader.GameLoader
// this function does not validate unit info and might panic if given bad unit info
func ModelToGameUnit(width, height uint8, unitInfo []byte) [][]objects.Unit {
	_units := make([][]objects.Unit, height)
	for i := uint8(0); i < height; i++ {
		_units[i] = make([]objects.Unit, width)
		for j := uint8(0); j < width; j++ {
			_units[i][j] = nil
		}
	}
	for i := 0; i < len(unitInfo); {
		y := unitInfo[i]
		x := unitInfo[i+1]
		p := uint8(unitInfo[i+2])
		t := unitInfo[i+3]
		hp := uint8(unitInfo[i+4])
		s := uint8(unitInfo[i+5])
		switch t {
		case objects.UnitTypeYou:
			_units[y][x] = objects.NewYou(p, hp, s)
			i += 6
		case objects.UnitTypeInfantry:
			_units[y][x] = objects.NewInfantry(p, hp, s)
			i += 6
		default:
			panic("panic convert: unknown unit type from unit info")
		}
	}
	return _units
}

// GameUnitToModel converts unit info from loader.GameLoader to model.Game
func GameUnitToModel(width, height uint8, _units [][]objects.Unit) []byte {
	var unitInfo []byte
	for i := uint8(0); i < height; i++ {
		for j := uint8(0); j < width; j++ {
			if _units[i][j] == nil {
				continue
			}

			unit := _units[i][j]
			switch unit.GetUnitType() {
			case objects.UnitTypeYou:
				you := unit.(*objects.You)
				unitInfo = append(unitInfo, i, j, you.Owner, unit.GetUnitType(), you.HP, unit.GetUnitState())
			case objects.UnitTypeInfantry:
				inf := unit.(*objects.Infantry)
				unitInfo = append(unitInfo, i, j, inf.Owner, unit.GetUnitType(), inf.HP, unit.GetUnitState())
			default:
				panic("panic convert: unknown unit type from unit object")
			}
		}
	}
	return unitInfo
}
