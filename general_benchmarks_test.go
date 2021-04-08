package beewarbe

import (
	"gitlab.com/beewar/beewar-be/configs"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"testing"
)

func BenchmarkGeneral(b *testing.B) {
	logger.InitLogger()
	defer logger.ShutdownLogger()
	configs.InitConfigs()
	access.InitAccess()
	defer access.ShutdownAccess()

	gameModel := access.QueryGameByID(1)
	gameUsers := access.QueryGameUsersByGameID(1)

	for i := 0; i < b.N; i++ {
		// do something
		gameModel.Status = 1
		_ = access.UpdateGameAndGameUser(gameModel, gameUsers)
	}
}
