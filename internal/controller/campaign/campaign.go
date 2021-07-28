package campaign

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

var campaignMapIDList []uint64

// InitCampaign initializes campaign module
func InitCampaign() {
	// fill campaign map id list
	maps, err := access.QueryCampaignMaps()
	if err != nil {
		logger.GetLogger().Fatal("error query campaign maps", zap.Error(err))
		return
	}

	campaignMapIDList = make([]uint64, len(maps)+1)
	for i, mapp := range maps {
		campaignMapIDList[i+1] = mapp.ID
	}

	logger.GetLogger().Info("loaded campaign maps", zap.Int("total", len(maps)))
}
