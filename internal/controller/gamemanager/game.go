package gamemanager

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/mapmanager"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// this file is the controller layer for http entry points

var (
	errMapDoesNotExist = errors.New("map does not exist")
	errMapNotReady     = errors.New("the map is not ready to be played. add more units")
)

// CreateGame creates a new game with the given map id. If password is provided, it will be bcrypt-ed
func CreateGame(mapID uint64, password string) (uint64, error) {
	mapModel := access.QueryMapByID(mapID)
	if mapModel == nil {
		return 0, errMapDoesNotExist
	}
	if !mapmanager.MapToGameValidation(mapModel) {
		return 0, errMapNotReady
	}
	var passwordHash = ""
	if password != "" {
		passwordHashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			logger.GetLogger().Error("error bcrypt", zap.Error(err))
			return 0, err
		}
		passwordHash = string(passwordHashByte)
	}
	return access.CreateGameFromMap(mapModel, passwordHash)
}

// GetMyGames returns a list of games for a specified user
func GetMyGames(userID uint64) ([]*model.GameUser, []*model.Game) {
	gameUsers := access.QueryGameUsersByUserID(userID)
	gameIDs := make([]uint64, len(gameUsers))
	for i := range gameUsers {
		gameIDs[i] = gameUsers[i].GameID
	}
	games := access.QueryGamesByID(gameIDs)
	return gameUsers, games
}
