package mapmanager

import (
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
)

// MapToGameValidation checks if a map is playable. Returns true if it is.
// This function does not do any validation on terrain/unit infos because they should be valid when saving the map.
func MapToGameValidation(mapModel *model.Map) bool {
	playerExist := make([]bool, mapModel.PlayerCount)
	playerExistCount := 0

	// this has some conversion logics too. Be careful when updating
	for i := 0; i < len(mapModel.UnitInfo); {
		p := int(mapModel.UnitInfo[i+2])
		t := mapModel.UnitInfo[i+3]
		switch t {
		case objects.UnitTypeYou:
			if playerExist[p-1] {
				return false
			}
			playerExist[p-1] = true
			playerExistCount++
			i += 6
		case objects.UnitTypeInfantry:
			i += 6
		default:
			panic("panic convert: unknown unit type from unit info")
		}
	}

	if playerExistCount != int(mapModel.PlayerCount) {
		return false
	}

	return true
}
