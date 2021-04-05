package gamemanager

import (
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
)

// mainly contains controllers for setup-ing game and deleting

var (
	errMapDoesNotExist  = errors.New("map does not exist")
	errGameDoesNotExist = errors.New("game does not exist")
)

// CreateGame creates a new game with the given map id
func CreateGame(mapID uint64) (uint64, error) {
	mapModel := access.QueryMapByID(mapID)
	if mapModel == nil {
		return 0, errMapDoesNotExist
	}
	return access.CreateGameFromMap(mapModel)
}

// RegisterForGame registers/links user to game
func RegisterForGame(userID, gameID uint64, playerOrder uint8) error {
	// we assume user to exist because it is provided from token
	if !access.IsExistGameByID(gameID) {
		return errGameDoesNotExist
	}
	// TODO: check if this user already registered to the game (user_id, game_id)
	// TODO: check if this playerOrder is already taken (game_id, player_order)
	if err := access.CreateGameUser(gameID, userID, playerOrder); err != nil {
		return err
	}
	return nil
}
