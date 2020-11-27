package accesstester

import (
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/logger"
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
	if int(mapp.Width)*int(mapp.Height) != 100 {
		return false
	}
	return true
}
