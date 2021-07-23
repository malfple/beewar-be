package seeder

import (
	"fmt"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/mapmanager"
)

// SeedSampleMaps inserts sample maps
func SeedSampleMaps() {
	var terrainInfo, unitInfo []byte

	// map 1
	if mapp, err := access.QueryMapByID(1); err == nil && mapp == nil {
		mapID, _ := mapmanager.CreateEmptyMap(1)
		fmt.Printf("create map with id: %d\n", mapID)
	}

	terrainInfo = []byte{
		1, 0, 1, 1, 3, 3, 1, 1, 0, 0,
		1, 0, 1, 1, 1, 3, 4, 1, 1, 0,
		0, 1, 1, 1, 1, 3, 3, 1, 1, 0,
		0, 1, 1, 1, 1, 2, 1, 1, 1, 1,
		1, 1, 1, 2, 0, 0, 2, 1, 1, 1,
		1, 1, 1, 2, 0, 0, 2, 1, 1, 1,
		1, 1, 1, 1, 2, 1, 1, 1, 1, 0,
		0, 1, 1, 5, 4, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 5, 4, 1, 1, 0, 1,
		0, 0, 1, 1, 5, 5, 1, 1, 0, 1,
	}
	unitInfo = []byte{
		5, 1, 1, 1, 10, 0,
		4, 1, 1, 3, 10, 0,
		3, 1, 1, 3, 10, 0,
		6, 1, 1, 3, 10, 0,
		7, 1, 1, 3, 10, 0,
		4, 8, 2, 1, 10, 0,
		3, 8, 2, 3, 10, 0,
		2, 8, 2, 3, 10, 0,
		5, 8, 2, 3, 10, 0,
		6, 8, 2, 3, 10, 0,
	}
	_ = access.UpdateMap(1, 0, 10, 10, "Test Map 1: Donut", 2,
		terrainInfo, unitInfo)
	fmt.Println("update map 1")

	// map 2
	if mapp, err := access.QueryMapByID(2); err == nil && mapp == nil {
		mapID, _ := mapmanager.CreateEmptyMap(1)
		fmt.Printf("create map with id: %d\n", mapID)
	}

	terrainInfo = []byte{
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
	}
	unitInfo = []byte{
		1, 7, 1, 1, 10, 0,
		1, 6, 1, 3, 10, 0,
		1, 8, 1, 3, 10, 0,
		0, 6, 1, 3, 10, 0,
		0, 7, 1, 3, 10, 0,
		2, 6, 1, 3, 10, 0,
		2, 7, 1, 3, 10, 0,
		4, 11, 2, 1, 10, 0,
		4, 10, 2, 3, 10, 0,
		4, 12, 2, 3, 10, 0,
		3, 11, 2, 3, 10, 0,
		3, 12, 2, 3, 10, 0,
		5, 11, 2, 3, 10, 0,
		5, 12, 2, 3, 10, 0,
		10, 2, 3, 1, 10, 0,
		10, 1, 3, 3, 10, 0,
		10, 3, 3, 3, 10, 0,
		9, 2, 3, 3, 10, 0,
		9, 3, 3, 3, 10, 0,
		11, 2, 3, 3, 10, 0,
		11, 3, 3, 3, 10, 0,
		13, 7, 4, 1, 10, 0,
		13, 6, 4, 3, 10, 0,
		13, 8, 4, 3, 10, 0,
		12, 6, 4, 3, 10, 0,
		12, 7, 4, 3, 10, 0,
		14, 6, 4, 3, 10, 0,
		14, 7, 4, 3, 10, 0,
	}
	_ = access.UpdateMap(2, 0, 15, 15, "Test Map 2: Cross", 4,
		terrainInfo, unitInfo)
	fmt.Println("update map 2")

	// map 3
	if mapp, err := access.QueryMapByID(3); err == nil && mapp == nil {
		mapID, _ := mapmanager.CreateEmptyMap(1)
		fmt.Printf("create map with id: %d\n", mapID)
	}

	terrainInfo = []byte{
		1, 1, 1, 1, 1, 1, 1, 5, 1, 4, 5, 3, 1, 5, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 3, 5, 4, 3, 3, 1, 3, 4, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 3, 1, 5, 5, 4, 5, 3, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 5, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	unitInfo = []byte{
		0, 0, 1, 1, 10, 0,
		0, 1, 1, 9, 4, 0,
		1, 1, 1, 9, 4, 0,
		2, 1, 1, 5, 10, 0,
		0, 2, 1, 3, 10, 0,
		1, 2, 1, 4, 8, 0,
		2, 2, 1, 6, 14, 0,
		0, 19, 2, 1, 10, 0,
		0, 18, 2, 9, 4, 0,
		1, 18, 2, 9, 4, 0,
		2, 18, 2, 5, 10, 0,
		0, 17, 2, 3, 10, 0,
		1, 17, 2, 4, 8, 0,
		2, 17, 2, 6, 14, 0,
	}
	_ = access.UpdateMap(3, 0, 4, 20, "Test Map 3: Line", 2,
		terrainInfo, unitInfo)
	fmt.Println("update map 3")
}
