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
	mapID, err := access.CreateEmptyMap(0, 2, 3, "some_map", 1, false)
	if err != nil {
		return false
	}
	mapModel, err := access.QueryMapByID(mapID)
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
	user1, err := access.QueryUserByUsername("player1")
	if err != nil || user1 == nil {
		return false
	}
	if err := access.CreateUser("player2@somemail.com", "player2", "password"); err != nil {
		return false
	}
	user2, err := access.QueryUserByUsername("player2")
	if err != nil || user2 == nil {
		return false
	}

	// game
	gameID, err := access.CreateGameFromMapUsingTx(nil, mapModel, "some game", "", 1)
	if err != nil {
		logger.GetLogger().Error("error create game from map", zap.Error(err))
		return false
	}
	for i, userID := range []uint64{user2.ID, user1.ID} {
		err = access.CreateGameUserUsingTx(nil, gameID, userID, uint8(i+1))
		if err != nil {
			logger.GetLogger().Error("error link game", zap.Error(err), zap.Uint64("user_id", userID))
			return false
		}
	}
	game, err := access.QueryGameByID(gameID)
	if err != nil {
		return false
	}
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

	games, err := access.QueryGamesByID([]uint64{gameID, 99999999})
	if err != nil {
		return false
	}
	if len(games) != 2 {
		logger.GetLogger().Error("games should have length 2")
		return false
	}
	if games[0].ID != gameID {
		logger.GetLogger().Error("games[0] id not match")
		return false
	}
	if games[1] != nil {
		logger.GetLogger().Error("games[1] is not nil")
		return false
	}

	// game users
	gus, err := access.QueryGameUsersByGameID(999999999)
	if err != nil {
		return false
	}
	if gus == nil {
		logger.GetLogger().Error("should be empty array, not nil")
		return false
	}
	gus, err = access.QueryGameUsersByUserID(999999999)
	if err != nil {
		return false
	}
	if gus == nil {
		logger.GetLogger().Error("should be empty array, not nil")
		return false
	}
	players, err := access.QueryGameUsersByGameID(gameID)
	if err != nil {
		return false
	}
	if len(players) != 2 || players[0].UserID != user2.ID || players[1].UserID != user1.ID {
		logger.GetLogger().Error("error query users linked to game")
		return false
	}
	games1, err := access.QueryGameUsersByUserID(user1.ID)
	if err != nil {
		return false
	}
	games2, err := access.QueryGameUsersByUserID(user1.ID)
	if err != nil {
		return false
	}
	if games1[0].GameID != gameID || games1[0].GameID != games2[0].GameID {
		logger.GetLogger().Error("error query games linked to user")
		return false
	}

	if gu, _ := access.QueryGameUser(gameID, user1.ID); gu == nil {
		logger.GetLogger().Error("game user should exist")
		return false
	}
	if gu, _ := access.QueryGameUser(gameID, 69696969); gu != nil {
		logger.GetLogger().Error("game user should not exist")
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
	gameUserToUpdate.MovesMade = 3
	if err := access.UpdateGameUserUsingTx(nil, gameUserToUpdate); err != nil {
		return false
	}
	games1again, _ := access.QueryGameUsersByUserID(user1.ID)
	if games1again[0].FinalTurns != 69 || games1again[0].FinalRank != 1 || games1again[0].MovesMade != 3 {
		logger.GetLogger().Error("mismatch update game user",
			zap.Int32("actual final turns", games1again[0].FinalTurns),
			zap.Uint8("actual final rank", games1again[0].FinalRank),
			zap.Uint32("actual moves made", games1again[0].MovesMade))
	}

	// combined updates
	game, _ = access.QueryGameByID(gameID)
	players, _ = access.QueryGameUsersByGameID(gameID)
	if err := access.UpdateGameAndGameUser(game, players); err != nil {
		return false
	}

	// queries
	_, err = access.QueryWaitingGames()
	if err != nil {
		return false
	}

	return true
}
