package mapmanager

import "errors"

const (
	// MapTypeFFA indicates a free-for-all map. Everyone is free to fight everyone else lol.
	MapTypeFFA = 0
	// MapTypeEscape indicates an escape mission map. 2 players only
	MapTypeEscape = 1
)

var (
	errMapTypeFFAPlayerCount = errors.New("FFA need 2+ players")
	errMapTypeInvalid = errors.New("invalid map type")
)

func validateMapType(mapType, playerCount uint8) error {
	switch mapType {
	case MapTypeFFA:
		if playerCount <= 1 {
			return errMapTypeFFAPlayerCount
		}
		return nil
	}
	return errMapTypeInvalid
}
