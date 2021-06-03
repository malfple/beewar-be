package beebot

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

const (
	errMsgBeebotUserMissing = "beebot user missing"
)

var beebotUser *model.User

// InitBeebotRoutines initializes bots. Only init after game manager is initialized.
// Does not need shutdown (it does automatically when gamemanager is shutdown)
func InitBeebotRoutines() {
	beebotUser = access.QueryUserByUsername("beebot")
	if beebotUser == nil {
		logger.GetLogger().Fatal("beebot user not found")
		return
	}

	gameUsers := access.QueryGameUsersByUserID(beebotUser.ID)
	for _, gu := range gameUsers {
		// small optimization to prevent starting up games where beebot already lost.
		if gu.FinalTurns != 0 {
			continue
		}
		// start for existing games.
		go startBeebotRoutine(gu.GameID, gu.PlayerOrder)
	}
}

// AskBeebotToJoinGame invites beebot to join a game
func AskBeebotToJoinGame(gameID uint64, playerOrder uint8, password string) string {
	if beebotUser == nil {
		return errMsgBeebotUserMissing
	}

	client := NewBotGameClient(beebotUser.ID, playerOrder)
	err := gamemanager.StartClientSession(client, gameID)
	if err != nil {
		return err.Error()
	}

	client.Hub.MessageBus <- &message.GameMessage{
		Cmd:    message.CmdJoin,
		Sender: beebotUser.ID,
		Data: &message.JoinMessageData{
			PlayerOrder: playerOrder,
			Password:    password,
		},
	}

	var reply *message.GameMessage
	for {
		reply = <-client.Replies
		if reply.Cmd == message.CmdJoin || reply.Cmd == message.CmdError {
			break
		}
	}

	gamemanager.EndClientSession(client)

	// haha. fail
	if reply.Cmd == message.CmdError {
		return reply.Data.(string)
	}

	// otherwise, start a goroutine
	go startBeebotRoutine(gameID, playerOrder)

	return ""
}

// This function setups the client and starts session.
func startBeebotRoutine(gameID uint64, playerOrder uint8) {
	if beebotUser == nil {
		return
	}
	client := NewBotGameClient(beebotUser.ID, playerOrder)
	err := gamemanager.StartClientSession(client, gameID)
	if err != nil {
		logger.GetLogger().Debug("error start beebot routine", zap.Error(err))
		return
	}
	logger.GetLogger().Debug("beebot client start listening", zap.Uint64("game_id", client.Hub.GameID))
	client.Listen()
	logger.GetLogger().Debug("beebot client stop listening", zap.Uint64("game_id", client.Hub.GameID))
	gamemanager.EndClientSession(client)
}
