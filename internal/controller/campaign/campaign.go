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

// 0-based. campaignMapList[0] is first campaign
var campaignMapList []*model.Map

// InitCampaign initializes campaign module
func InitCampaign() {
	// fill campaign map id list
	maps, err := access.QueryCampaignMaps()
	if err != nil {
		logger.GetLogger().Fatal("error query campaign maps", zap.Error(err))
		return
	}

	campaignMapList = maps

	logger.GetLogger().Info("loaded campaign maps", zap.Int("total", len(maps)))
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
	if int32(campaignLevel) > user.HighestCampaign+1 {
		return 0, errCampaignNotAccessible
	}
	if user.CurrCampaignGameID != 0 {
		// check that game's status
		game, err := access.QueryGameByID(user.CurrCampaignGameID)
		if err != nil {
			return 0, err
		}
		if game.Status != loader.GameStatusEnded {
			return 0, errOngoingCampaign
		}
		// game already ended, it's ok
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
