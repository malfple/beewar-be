package seeder

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
)

func seedCampaignMap2(existingMaps []*model.Map) bool {
	// if existing map id is 0, create new
	mapID := uint64(0)
	if 1 < len(existingMaps) {
		mapID = existingMaps[1].ID
	}

	if mapID == 0 {
		id, err := access.CreateEmptyMap(0, 8, 10, "Campaign Map 2", 1, true)
		if err != nil {
			return false
		}
		mapID = id
		fmt.Println("create new campaign map 2, id: ", mapID)
	}

	terrainInfo := []byte{
		0, 5, 5, 5, 0, 0, 0, 0, 0, 0,
		4, 4, 4, 5, 5, 5, 0, 0, 0, 0,
		4, 5, 5, 5, 2, 2, 2, 2, 2, 2,
		4, 5, 4, 5, 2, 2, 3, 1, 1, 1,
		5, 4, 5, 2, 2, 3, 2, 1, 1, 1,
		1, 1, 1, 5, 2, 3, 2, 2, 1, 1,
		1, 1, 1, 1, 1, 2, 2, 2, 1, 1,
		0, 5, 2, 2, 2, 2, 2, 2, 2, 1,
	}
	unitInfo := []byte{
		4, 5, 1, 3, 10, 0,
		4, 7, 2, 3, 10, 0,
		5, 5, 1, 3, 10, 0,
		5, 8, 2, 3, 10, 0,
		6, 2, 1, 1, 10, 0,
		6, 3, 1, 3, 10, 0,
		6, 8, 2, 3, 10, 0,
		6, 9, 2, 3, 10, 0,
		7, 9, 2, 1, 10, 0,
	}
	if err := access.UpdateMap(mapID, 0, 8, 10, "Campaign Map 2", 2,
		terrainInfo, unitInfo); err != nil {
		return false
	}
	fmt.Println("update campaign map 2, id: ", mapID)

	return true
}
