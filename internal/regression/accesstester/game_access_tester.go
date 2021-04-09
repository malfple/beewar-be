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
	mapModel := access.QueryMapByID(mapID)

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
	gameID, err := access.CreateGameFromMap(mapModel, "")
	if err != nil {
		logger.GetLogger().Error("error create game from map", zap.Error(err))
		return false
	}
	for i, userID := range []uint64{user2.ID, user1.ID} {
		err = access.CreateGameUser(gameID, userID, uint8(i+1))
		if err != nil {
			logger.GetLogger().Error("error link game", zap.Error(err), zap.Uint64("user_id", userID))
			return false
		}
	}
	game := access.QueryGameByID(gameID)
	if game == nil {
		logger.GetLogger().Error("game not found")
		return false
	}
	game.TurnPlayer = 2
	game.TurnCount = 1
	game.UnitInfo = append(game.UnitInfo, []byte{1, 2, 1, 1, 10, 0}...)
	err = access.UpdateGameUsingTx(nil, game)
	if err != nil {
		logger.GetLogger().Error("error update game", zap.Error(err))
		return false
	}

	// game users
	if access.QueryGameUsersByGameID(999999999) == nil {
		logger.GetLogger().Error("should be empty array, not nil")
	}
	if access.QueryGameUsersByUserID(999999999) == nil {
		logger.GetLogger().Error("should be empty array, not nil")
	}
	players := access.QueryGameUsersByGameID(gameID)
	if len(players) != 2 || players[0].UserID != user2.ID || players[1].UserID != user1.ID {
		logger.GetLogger().Error("error query users linked to game")
		return false
	}
	games1 := access.QueryGameUsersByUserID(user1.ID)
	games2 := access.QueryGameUsersByUserID(user1.ID)
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

	if !access.IsExistGameUserByLink(user1.ID, gameID) {
		logger.GetLogger().Error("game user not found")
		return false
	}
	if access.IsExistGameUserByLink(696969, gameID) {
		logger.GetLogger().Error("game user isn't supposed to exist")
		return false
	}
	if !access.IsExistGameUserByPlayerOrder(gameID, 2) {
		logger.GetLogger().Error("slot 2 taken")
		return false
	}
	if access.IsExistGameUserByPlayerOrder(gameID, 3) {
		logger.GetLogger().Error("slot 3 doesn't exist haha")
		return false
	}

	// update game user
	gameUserToUpdate := games1[0]
	gameUserToUpdate.FinalTurns = 69
	gameUserToUpdate.FinalRank = 1
	if err := access.UpdateGameUserUsingTx(nil, gameUserToUpdate); err != nil {
		return false
	}
	games1again := access.QueryGameUsersByUserID(user1.ID)
	if games1again[0].FinalTurns != 69 || games1again[0].FinalRank != 1 {
		logger.GetLogger().Error("mismatch update game user",
			zap.Int32("actual final turns", games1again[0].FinalTurns),
			zap.Uint8("actual final rank", games1again[0].FinalRank))
	}

	// combined updates
	game = access.QueryGameByID(gameID)
	players = access.QueryGameUsersByGameID(gameID)
	if err := access.UpdateGameAndGameUser(game, players); err != nil {
		return false
	}

	// queries
	_ = access.QueryWaitingGames()

	return true
}
