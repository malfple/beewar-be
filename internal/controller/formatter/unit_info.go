package formatter

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
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
p = the player who owns this unit (0..number of players, 0 means neutral, >=1 means player owned)
t = unit type
hp = unit health point
s = state/flags (unit state)

no two units can share the same position

the most basic information that a flag has are:
- turn state - most units only have two states: not yet moved, and already moved: hence 1 bit
	- this complicates more when the units store firing states and can carry other units


unit type information is available in `units` package
*/

var (
	// ErrMapInvalidUnitInfo is returned when unit info does not follow format
	ErrMapInvalidUnitInfo = errors.New("invalid unit info")
	// ErrMapUnitSamePosition is returned when two units are in the same position
	ErrMapUnitSamePosition = errors.New("no two units can share the same position")
)

// ValidateUnitInfo validates whether unit info follows format
func ValidateUnitInfo(height, width int, unitInfo []byte) error {
	posMap := make(map[int]bool)
	for i := 0; i < len(unitInfo); {
		if i+5 >= len(unitInfo) { // the remaining length is less than 6 (the required minimum of a normal unit)
			return ErrMapInvalidUnitInfo
		}
		y := int(unitInfo[i])
		x := int(unitInfo[i+1])
		if y < 0 || y >= height || x < 0 || x >= width {
			return ErrMapInvalidUnitInfo
		}
		if _, ok := posMap[y*width+x]; ok { // 2 unit in the same pos
			return ErrMapUnitSamePosition
		}
		posMap[y*width+x] = true
		t := unitInfo[i+3]
		switch t {
		case objects.UnitTypeQueen:
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
func ModelToGameUnit(height, width int, unitInfo []byte) [][]objects.Unit {
	_units := make([][]objects.Unit, height)
	for i := 0; i < height; i++ {
		_units[i] = make([]objects.Unit, width)
		for j := 0; j < width; j++ {
			_units[i][j] = nil
		}
	}
	for i := 0; i < len(unitInfo); {
		y := unitInfo[i]
		x := unitInfo[i+1]
		p := int(unitInfo[i+2])
		t := unitInfo[i+3]
		hp := int(unitInfo[i+4])
		s := int(unitInfo[i+5])
		switch t {
		case objects.UnitTypeQueen:
			_units[y][x] = objects.NewQueen(p, hp, s)
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
func GameUnitToModel(height, width int, _units [][]objects.Unit) []byte {
	var unitInfo []byte
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if _units[i][j] == nil {
				continue
			}

			unit := _units[i][j]
			switch unit.UnitType() {
			case objects.UnitTypeQueen:
				queen := unit.(*objects.Queen)
				unitInfo = append(unitInfo, byte(i), byte(j), byte(queen.Owner), byte(unit.UnitType()), byte(queen.HP), byte(unit.GetState()))
			case objects.UnitTypeInfantry:
				inf := unit.(*objects.Infantry)
				unitInfo = append(unitInfo, byte(i), byte(j), byte(inf.Owner), byte(unit.UnitType()), byte(inf.HP), byte(unit.GetState()))
			default:
				panic("panic convert: unknown unit type from unit object")
			}
		}
	}
	return unitInfo
}
