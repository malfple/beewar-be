package gamemanager

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// mainly contains controllers for setup-ing game and deleting

var (
	errMapDoesNotExist   = errors.New("map does not exist")
	errGameDoesNotExist  = errors.New("game does not exist")
	errAlreadyRegistered = errors.New("user already registered for this game")
	errPlayerOrderTaken  = errors.New("that slot/player_order is already taken")
)

// CreateGame creates a new game with the given map id. If password is provided, it will be bcrypt-ed
func CreateGame(mapID uint64, password string) (uint64, error) {
	mapModel := access.QueryMapByID(mapID)
	if mapModel == nil {
		return 0, errMapDoesNotExist
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

// RegisterForGame registers/links user to game
func RegisterForGame(userID, gameID uint64, playerOrder uint8) error {
	// we assume user to exist because it is provided from token
	if !access.IsExistGameByID(gameID) {
		return errGameDoesNotExist
	}
	if access.IsExistGameUserByLink(userID, gameID) {
		return errAlreadyRegistered
	}
	if access.IsExistGameUserByPlayerOrder(gameID, playerOrder) {
		return errPlayerOrderTaken
	}
	if err := access.CreateGameUser(gameID, userID, playerOrder); err != nil {
		return err
	}
	return nil
}
