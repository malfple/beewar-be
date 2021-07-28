package seeder

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// SeedCampaignMaps inserts campaign maps
func SeedCampaignMaps() bool {
	existingMaps, err := access.QueryCampaignMaps()
	if err != nil {
		logger.GetLogger().Error("error query existing campaign maps", zap.Error(err))
		return false
	}

	if !seedCampaignMap1(existingMaps) {
		return false
	}
	return true
}
