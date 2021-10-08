package campaign

import (
	"database/sql"
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/beebot"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/loader"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

var (
	errCampaignNotExist      = errors.New("that campaign does not exist (yet)")
	errCampaignNotAccessible = errors.New("that campaign is not yet accessible by you")
	errOngoingCampaign       = errors.New("there is already an ongoing campaign")
)

// Campaign map level is determined by the order in the db (sorted by id)
// The seeder is responsible for keeping this order

// 0-based. campaignMapList[0] is first campaign
var campaignMapList []*model.Map
var campaignMapIDToLevelMap map[uint64]int

// InitCampaign initializes campaign module
func InitCampaign() {
	// fill campaign map id list
	maps, err := access.QueryCampaignMaps()
	if err != nil {
		logger.GetLogger().Fatal("error query campaign maps", zap.Error(err))
		return
	}

	campaignMapList = maps

	// create id to level map
	campaignMapIDToLevelMap = make(map[uint64]int)
	for i, campaignMap := range campaignMapList {
		campaignMapIDToLevelMap[campaignMap.ID] = i + 1
	}

	logger.GetLogger().Info("loaded campaign maps", zap.Int("total", len(maps)))
}

// GetCampaignList returns the list of campaign maps
func GetCampaignList() []*model.Map {
	return campaignMapList
}

// Checks current campaign, and update user if possible (if game already ended).
// This function does not touch db, only updates the user model.
// Then it returns the game if it exists and is ongoing.
func getCurrentCampaignAndUpdateUser(user *model.User) (*model.Game, error) {
	if user.CurrCampaignGameID == 0 {
		return nil, nil
	}

	game, err := access.QueryGameByID(user.CurrCampaignGameID)
	if err != nil {
		return nil, err
	}
	if game.Status != loader.GameStatusEnded {
		return game, nil // game found
	}
	// game already ended, update user
	// check if the user actually won
	gameUser, err := access.QueryGameUser(game.ID, user.ID)
	if err != nil {
		return nil, err
	}
	if gameUser.FinalRank == 1 { // user won
		if level, ok := campaignMapIDToLevelMap[game.MapID]; ok {
			if level > int(user.HighestCampaign) {
				user.HighestCampaign = int32(level)
			}
		}
	}
	user.CurrCampaignGameID = 0
	// does not touch db
	return nil, nil
}

// StartNewCampaign starts a new campaign for the user.
// Also returns error if the user currently has an ongoing campaign.
// Returns (game_id, error)
func StartNewCampaign(userID uint64, campaignLevel int) (uint64, error) {
	user, err := access.QueryUserByID(userID)
	if err != nil {
		return 0, err
	}
	if campaignLevel < 1 || campaignLevel > len(campaignMapList) {
		return 0, errCampaignNotExist
	}
	{
		game, err := getCurrentCampaignAndUpdateUser(user)
		if err != nil {
			return 0, nil
		}
		if game != nil {
			return 0, errOngoingCampaign
		}
	}
	// this check has to be after the current campaign check to possibly update user.HighestCampaign
	if int32(campaignLevel) > user.HighestCampaign+1 {
		return 0, errCampaignNotAccessible
	}

	campaignMap := campaignMapList[campaignLevel-1]

	var gameID uint64
	// start new game here, prefilled with the players
	if err := access.ExecWithTransaction(func(tx *sql.Tx) error {
		var err error
		gameID, err = access.CreateGameFromMapUsingTx(tx, campaignMap,
			campaignMap.Name, "", user.ID)
		if err != nil {
			return err
		}
		if err := access.CreateGameUserUsingTx(tx, gameID, user.ID, 1); err != nil {
			return err
		}
		if err := access.CreateGameUserUsingTx(tx, gameID, beebot.GetBeebotUserID(), 2); err != nil {
			return err
		}
		// assign this game as current campaign game
		user.CurrCampaignGameID = gameID
		if err := access.UpdateUserUsingTx(tx, user); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	// if success, start a beebot routine
	go beebot.RunBeebotRoutine(gameID, 2)

	return gameID, nil
}

// GetCurrentCampaign returns the current campaign game (game_id, campaign_level) by a user.
// returns game_id = 0 if no game found.
func GetCurrentCampaign(userID uint64) (uint64, int32, error) {
	user, err := access.QueryUserByID(userID)
	if err != nil {
		return 0, 0, err
	}
	game, err := getCurrentCampaignAndUpdateUser(user)
	if err != nil {
		return 0, 0, err
	}
	if game != nil {
		return game.ID, int32(campaignMapIDToLevelMap[game.MapID]), nil // game exists
	}
	_ = access.UpdateUserUsingTx(nil, user) // fail silently
	return 0, 0, nil
}
