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

	mapID, err := access.CreateEmptyMap(0, 2, 3, "some_map", 1)
	if err != nil {
		return false
	}
	_, err = access.CreateEmptyMap(0, 3, 2, "map2", 1)
	if err != nil {
		return false
	}
	_, err = access.CreateEmptyMap(0, 5, 5, "map3", 1)
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

	mapp := access.QueryMapByID(mapID)
	if mapp == nil {
		return false
	}
	if mapp.ID != mapID {
		return false
	}
	if mapp.Width*mapp.Height != 100 {
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
	if access.QueryMaps(0, 0) == nil {
		logger.GetLogger().Error("maps should be empty array, not nil")
		return false
	}

	maps := access.QueryMaps(2, 0)
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

	maps = access.QueryMaps(10, 1)
	if maps[0].Name != "map2" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "map2"), zap.String("found", maps[1].Name))
	}
	if maps[1].Name != "some updated map" {
		logger.GetLogger().Error("mismatch map name", zap.String("expected", "aome updated map"), zap.String("found", maps[1].Name))
	}

	maps = access.QueryMaps(10, 10)
	if len(maps) != 0 {
		logger.GetLogger().Error("mismatch number of maps", zap.Int("expected", 0), zap.Int("found", len(maps)))
		return false
	}

	return true
}
