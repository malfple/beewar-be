package seeder

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
)

func seedCampaignMap1(existingMaps []*model.Map) bool {
	// if existing map id is 0, create new
	mapID := uint64(0)
	if 0 < len(existingMaps) {
		mapID = existingMaps[0].ID
	}

	if mapID == 0 {
		id, err := access.CreateEmptyMap(0, 10, 10, "Campaign Map 1", 1, true)
		if err != nil {
			return false
		}
		mapID = id
		fmt.Println("create new campaign map 1, id: ", mapID)
	}

	terrainInfo := []byte{
		0, 0, 0, 1, 0, 0, 0, 0, 1, 0,
		0, 1, 1, 0, 0, 1, 1, 1, 0, 1,
		0, 1, 1, 0, 0, 1, 1, 1, 0, 0,
		0, 2, 0, 1, 1, 2, 1, 1, 0, 2,
		1, 0, 1, 1, 1, 1, 1, 0, 1, 0,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 1, 1, 1, 1, 0, 1, 2,
		0, 0, 0, 1, 1, 2, 1, 0, 1, 0,
		2, 1, 1, 0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 1, 1, 0, 1, 1, 2,
	}
	unitInfo := []byte{
		1, 7, 2, 3, 10, 0,
		2, 7, 2, 1, 10, 0,
		5, 3, 1, 3, 10, 0,
		6, 3, 1, 1, 10, 0,
		7, 4, 1, 3, 10, 0,
	}
	if err := access.UpdateMap(mapID, 0, 10, 10, "Campaign Map 1", 2,
		terrainInfo, unitInfo); err != nil {
		return false
	}
	fmt.Println("update campaign map 1, id: ", mapID)

	return true
}
