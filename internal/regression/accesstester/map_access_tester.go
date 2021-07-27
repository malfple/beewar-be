package accesstester

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"math/rand"
)

// TestMapAccess runs regression tests for map access
func TestMapAccess() bool {
	logger.GetLogger().Info("map access tester")

	mapID, err := access.CreateEmptyMap(0, 2, 3, "some_map", 1, false)
	if err != nil {
		return false
	}
	mapID2, err := access.CreateEmptyMap(0, 3, 2, "map2", 1, true)
	if err != nil {
		return false
	}
	_, err = access.CreateEmptyMap(0, 5, 5, "map3", 1, true)
	if err != nil {
		return false
	}

	terrain1 := make([]byte, 100)
	for i := 0; i < 10; i++ {
		terrain1[rand.Int()%100] = 1
	}
	err = access.UpdateMap(mapID, 0, 10, 10, "some updated map", 2, terrain1, make([]byte, 0))
	if err != nil {
		logger.GetLogger().Error("error update map", zap.Error(err))
		return false
	}

	mapp, err := access.QueryMapByID(mapID)
	if err != nil {
		return false
	}
	if mapp == nil {
		return false
	}
	if mapp.ID != mapID {
		return false
	}
	if mapp.Width*mapp.Height != 100 {
		return false
	}
	if mapp.IsCampaign != false {
		return false
	}

	mapp2, err := access.QueryMapByID(mapID2)
	if err != nil {
		return false
	}
	if mapp2.IsCampaign != true {
		return false
	}

	// batch queries
	if !TestMapAccessQueryMaps() {
		return false
	}

	return true
}

// TestMapAccessQueryMaps tests access.QueryMaps
func TestMapAccessQueryMaps() bool {
	// query 0 maps, should not be null
	if maps, _ := access.QueryMaps(0, 0); maps == nil {
		logger.GetLogger().Error("maps should be empty array, not nil")
		return false
	}

	maps, err := access.QueryMaps(2, 0)
	if err != nil {
		return false
	}
	if len(maps) != 2 {
		logger.GetLogger().Error("mismatch number of maps", zap.Int("expected", 2), zap.Int("found", len(maps)))
		return false
	}
	if maps[0].Name != "map3" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "map3"), zap.String("found", maps[0].Name))
	}
	if maps[1].Name != "map2" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "map2"), zap.String("found", maps[1].Name))
	}

	maps, err = access.QueryMaps(10, 1)
	if err != nil {
		return false
	}
	if maps[0].Name != "map2" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "map2"), zap.String("found", maps[1].Name))
	}
	if maps[1].Name != "some updated map" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "aome updated map"), zap.String("found", maps[1].Name))
	}

	maps, err = access.QueryMaps(10, 10)
	if err != nil {
		return false
	}
	if len(maps) != 0 {
		logger.GetLogger().Error("mismatch number of maps", zap.Int("expected", 0), zap.Int("found", len(maps)))
		return false
	}

	// query campaign maps -> only made 1
	maps, err = access.QueryCampaignMaps()
	if err != nil {
		return false
	}
	if len(maps) != 2 {
		logger.GetLogger().Error("mismatch number of maps", zap.Int("expected", 2), zap.Int("found", len(maps)))
		return false
	}
	if maps[0].ID > maps[1].ID {
		logger.GetLogger().Error("campaign maps not ordered")
		return false
	}

	return true
}
