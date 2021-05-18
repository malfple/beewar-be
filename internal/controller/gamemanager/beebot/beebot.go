package beebot

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

const (
	errMsgBeebotUserMissing = "beebot user missing"
)

// InitBeebotRoutines initializes bots. Only init after game manager is initialized.
// Does not need shutdown (it does automatically when gamemanager is shutdown)
func InitBeebotRoutines() {
	beebotUser := access.QueryUserByUsername("beebot")
	if beebotUser == nil {
		logger.GetLogger().Fatal("beebot user not found")
		return
	}

	gameUsers := access.QueryGameUsersByUserID(beebotUser.ID)
	for _, gu := range gameUsers {
		// start for existing games. game state doesn't matter, it will auto-close when necessary anyway.
		go startBeebotRoutine(beebotUser.ID, gu.GameID)
	}
}

// AskBeebotToJoinGame invites beebot to join a game
func AskBeebotToJoinGame(gameID uint64, playerOrder uint8, password string) string {
	beebotUser := access.QueryUserByUsername("beebot")
	if beebotUser == nil {
		return errMsgBeebotUserMissing
	}

	client := NewBotGameClient(beebotUser.ID)
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
	go startBeebotRoutine(beebotUser.ID, gameID)

	return ""
}

// This function setups the client and starts session.
func startBeebotRoutine(botUserID, gameID uint64) {
	client := NewBotGameClient(botUserID)
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
