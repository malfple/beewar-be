package loader

import (
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"golang.org/x/crypto/bcrypt"
)

// handles player join
func (gl *GameLoader) handleJoin(msg *message.GameMessage) (*message.GameMessage, bool) {
	if gl.Status != GameStatusPicking {
		return message.GameErrorMessage(errMsgGameNotInPicking), false
	}
	data := msg.Data.(*message.JoinMessageData)
	if data.PlayerOrder < 1 || int(data.PlayerOrder) > gl.PlayerCount {
		return message.GameErrorMessage(errMsgPlayerOrderInvalid), false
	}
	if gl.GameUsers[data.PlayerOrder-1].UserID != 0 { // taken haha
		return message.GameErrorMessage(errMsgPlayerOrderTaken), false
	}
	if _, ok := gl.UserIDToPlayerMap[msg.Sender]; ok {
		return message.GameErrorMessage(errMsgAlreadyJoined), false
	}
	if gl.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(gl.Password), []byte(data.Password)); err != nil {
			return message.GameErrorMessage(errMsgWrongPassword), false
		}
	}
	// pass join validation
	if err := access.CreateGameUser(gl.ID, msg.Sender, data.PlayerOrder); err != nil {
		return message.GameErrorMessage(err.Error()), false
	}
	gl.GameUsers[data.PlayerOrder-1], _ = access.QueryGameUser(gl.ID, msg.Sender) // TODO: handle error
	users, _ := access.QueryUsersByID([]uint64{msg.Sender})                       // TODO: handle error
	gl.Users[data.PlayerOrder-1] = users[0]
	gl.Users[data.PlayerOrder-1].Password = "nope..."
	gl.UserIDToPlayerMap[msg.Sender] = int(data.PlayerOrder)
	gl.checkGameStart()
	return &message.GameMessage{
		Cmd:    msg.Cmd,
		Sender: msg.Sender,
		Data: &message.JoinMessageDataExt{
			Player: &message.Player{
				UserID:      msg.Sender,
				PlayerOrder: data.PlayerOrder,
				FinalRank:   0,
				FinalTurns:  0,
				User:        gl.Users[data.PlayerOrder-1],
			},
		},
	}, true
}
