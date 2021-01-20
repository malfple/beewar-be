package accesstester

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"math/rand"
)

// TestGameAccess runs regression tests for game access and game user access
func TestGameAccess() bool {
	logger.GetLogger().Info("game access tester")

	// create map first
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

	// create players
	if err := access.CreateUser("player1@somemail.com", "player1", "password"); err != nil {
		return false
	}
	user1 := access.QueryUserByUsername("player1")
	if user1 == nil {
		return false
	}
	if err := access.CreateUser("player2@somemail.com", "player2", "password"); err != nil {
		return false
	}
	user2 := access.QueryUserByUsername("player2")
	if user2 == nil {
		return false
	}

	// game
	gameID, err := access.CreateGameFromMap(mapID, []uint64{user1.ID, user2.ID})
	if err != nil {
		logger.GetLogger().Error("error create game from map", zap.Error(err))
		return false
	}
	game := access.QueryGameByID(gameID)
	if game == nil {
		logger.GetLogger().Error("game not found")
		return false
	}
	game.TurnPlayer = 2
	game.TurnCount = 1
	game.UnitInfo = append(game.UnitInfo, []byte{1, 2, 1, 1, 10, 0}...)
	err = access.UpdateGame(game)
	if err != nil {
		logger.GetLogger().Error("error update game", zap.Error(err))
		return false
	}

	players := access.QueryUsersLinkedToGame(gameID)
	if len(players) != 2 || players[0].UserID != user1.ID || players[1].UserID != user2.ID {
		logger.GetLogger().Error("error query users linked to game")
		return false
	}
	games1 := access.QueryGamesLinkedToUser(user1.ID)
	games2 := access.QueryGamesLinkedToUser(user1.ID)
	if games1[0].GameID != gameID || games1[0].GameID != games2[0].GameID {
		logger.GetLogger().Error("error query games linked to user")
		return false
	}

	if !access.IsExistGameByID(gameID) {
		logger.GetLogger().Error("game doesn't exist")
		return false
	}
	if access.IsExistGameByID(696969) {
		logger.GetLogger().Error("game isn't supposed to exist")
		return false
	}

	return true
}
