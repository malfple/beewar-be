package mapmanager

import "errors"

const (
	// MapTypeFFA indicates a free-for-all map. Everyone is free to fight everyone else lol.
	MapTypeFFA = 0
	// MapTypeEscape indicates an escape mission map. 2 players only
	MapTypeEscape = 1
)

var (
	errMapTypeFFAPlayerCount    = errors.New("FFA need 2+ players")
	errMapTypeEscapePlayerCount = errors.New("escape mission needs exactly 2 players")
	errMapTypeInvalid           = errors.New("invalid map type")
)

func validateMapType(mapType, playerCount uint8) error {
	switch mapType {
	case MapTypeFFA:
		if playerCount <= 1 {
			return errMapTypeFFAPlayerCount
		}
		return nil
	case MapTypeEscape:
		if playerCount != 2 {
			return errMapTypeEscapePlayerCount
		}
		return nil
	}
	return errMapTypeInvalid
}

// returns -1 if there is no limit
func calcExpectedThroneCount(mapType uint8) int {
	switch mapType {
	case MapTypeEscape:
		return 1
	default:
		return -1
	}
}
