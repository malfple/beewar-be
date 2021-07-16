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
	// ErrMapUnitOutsideMap is returned when some units are outside map
	ErrMapUnitOutsideMap = errors.New("unit outside map")
	// ErrMapUnitSamePosition is returned when two units are in the same position
	ErrMapUnitSamePosition = errors.New("no two units can share the same position")
	// ErrMapPlayerNotExist is returned when a unit is owned by a player with number greater than the player count
	ErrMapPlayerNotExist = errors.New("some units belong to non-existent player. raise player count to fix this")
	// ErrMapPlayerQueen is returned when some queens are missing or duplicated
	ErrMapPlayerQueen = errors.New("queen count has to match player count, some player have missing or duplicated queen")
)

// ValidateUnitInfo validates whether unit info follows format.
// the `skipGameReadiness` param should be set to true when saving a game unit info, to disable validations for game-readiness
// (like queen count and duplication)
func ValidateUnitInfo(height, width, playerCount int, unitInfo []byte, skipGameReadiness bool) error {
	queenExist := make([]bool, playerCount)
	queenCount := 0
	posMap := make(map[int]bool)
	for i := 0; i < len(unitInfo); {
		if i+5 >= len(unitInfo) { // the remaining length is less than 6 (the required minimum of a normal unit)
			return ErrMapInvalidUnitInfo
		}
		y := int(unitInfo[i])
		x := int(unitInfo[i+1])
		if y < 0 || y >= height || x < 0 || x >= width {
			return ErrMapUnitOutsideMap
		}
		if _, ok := posMap[y*width+x]; ok { // 2 unit in the same pos
			return ErrMapUnitSamePosition
		}
		posMap[y*width+x] = true
		p := unitInfo[i+2]
		t := unitInfo[i+3]
		if int(p) > playerCount {
			return ErrMapPlayerNotExist
		}
		switch t {
		case objects.UnitTypeQueen:
			if !skipGameReadiness {
				if queenExist[p-1] {
					return ErrMapPlayerQueen
				}
				queenExist[p-1] = true
				queenCount++
			}
			i += 6
		case objects.UnitTypeInfantry:
			i += 6
		case objects.UnitTypeJetCrew:
			i += 6
		case objects.UnitTypeWizard:
			i += 6
		case objects.UnitTypeTank:
			i += 6
		case objects.UnitTypeMortar:
			i += 6
		default:
			return ErrMapInvalidUnitInfo
		}
	}
	if !skipGameReadiness && playerCount != queenCount {
		return ErrMapPlayerQueen
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
		case objects.UnitTypeJetCrew:
			_units[y][x] = objects.NewJetCrew(p, hp, s)
			i += 6
		case objects.UnitTypeWizard:
			_units[y][x] = objects.NewWizard(p, hp, s)
			i += 6
		case objects.UnitTypeTank:
			_units[y][x] = objects.NewTank(p, hp, s)
			i += 6
		case objects.UnitTypeMortar:
			_units[y][x] = objects.NewMortar(p, hp, s)
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
			// most units have similar format
			case objects.UnitTypeQueen:
				fallthrough
			case objects.UnitTypeInfantry:
				fallthrough
			case objects.UnitTypeJetCrew:
				fallthrough
			case objects.UnitTypeWizard:
				fallthrough
			case objects.UnitTypeTank:
				fallthrough
			case objects.UnitTypeMortar:
				unitInfo = append(unitInfo, byte(i), byte(j),
					byte(unit.GetOwner()), byte(unit.UnitType()), byte(unit.GetHP()), byte(unit.GetState()))
			default:
				panic("panic convert: unknown unit type from unit object")
			}
		}
	}
	return unitInfo
}
