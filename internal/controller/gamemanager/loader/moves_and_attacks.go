package loader

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/cmdwhitelist"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/combat"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
)

func (gl *GameLoader) handleUnitMove(msg *message.GameMessage) (*message.GameMessage, bool) {
	data := msg.Data.(*message.UnitMoveMessageData)
	if errMsg := gl.validateUnitOwned(msg.Sender, data.Y1, data.X1); len(errMsg) > 0 {
		return message.GameErrorMessage(errMsg), false
	}
	if _, ok := cmdwhitelist.UnitMoveMap[gl.Units[data.Y1][data.X1].GetUnitType()]; !ok {
		return message.GameErrorMessage(ErrMsgUnitCmdNotAllowed), false
	}
	if gl.Units[data.Y1][data.X1].GetUnitStateBit(objects.UnitStateBitMoved) { // has this unit moved?
		return message.GameErrorMessage(ErrMsgUnitAlreadyMoved), false
	}
	if !gl.GridEngine.ValidateMove(data.Y1, data.X1, data.Y2, data.X2) {
		return message.GameErrorMessage(ErrMsgInvalidMove), false
	}
	gl.Units[data.Y2][data.X2], gl.Units[data.Y1][data.X1] = gl.Units[data.Y1][data.X1], gl.Units[data.Y2][data.X2]
	gl.Units[data.Y2][data.X2].ToggleUnitStateBit(objects.UnitStateBitMoved)
	return msg, true
}

func (gl *GameLoader) handleUnitAttack(msg *message.GameMessage) (*message.GameMessage, bool) {
	data := msg.Data.(*message.UnitAttackMessageData)
	if errMsg := gl.validateUnitOwned(msg.Sender, data.Y1, data.X1); len(errMsg) > 0 {
		return message.GameErrorMessage(errMsg), false
	}
	if _, ok := cmdwhitelist.UnitAttackMap[gl.Units[data.Y1][data.X1].GetUnitType()]; !ok {
		return message.GameErrorMessage(ErrMsgUnitCmdNotAllowed), false
	}
	if gl.Units[data.Y1][data.X1].GetUnitStateBit(objects.UnitStateBitMoved) { // has this unit moved?
		return message.GameErrorMessage(ErrMsgUnitAlreadyMoved), false
	}
	okAtk, distAtk := gl.GridEngine.ValidateAttack(data.Y1, data.X1, data.YT, data.XT, gl.Units[data.Y1][data.X1])
	if !okAtk {
		return message.GameErrorMessage(ErrMsgInvalidAttack), false
	}
	gl.Units[data.Y1][data.X1].ToggleUnitStateBit(objects.UnitStateBitMoved)
	combat.NormalCombat(gl.Units[data.Y1][data.X1], gl.Units[data.YT][data.XT], distAtk)
	replyMsg := &message.GameMessage{
		Cmd:    msg.Cmd,
		Sender: msg.Sender,
		Data: &message.UnitAttackMessageDataExt{
			Y1:    data.Y1,
			X1:    data.X1,
			YT:    data.YT,
			XT:    data.XT,
			HPAtk: gl.Units[data.Y1][data.X1].GetUnitHP(),
			HPDef: gl.Units[data.YT][data.XT].GetUnitHP(),
		},
	}
	gl.checkUnitAlive(data.Y1, data.X1)
	gl.checkUnitAlive(data.YT, data.XT)
	return replyMsg, true
}

func (gl *GameLoader) handleUnitMoveAndAttack(msg *message.GameMessage) (*message.GameMessage, bool) {
	data := msg.Data.(*message.UnitMoveAndAttackMessageData)
	if errMsg := gl.validateUnitOwned(msg.Sender, data.Y1, data.X1); len(errMsg) > 0 {
		return message.GameErrorMessage(errMsg), false
	}
	if _, ok := cmdwhitelist.UnitMoveAndAttackMap[gl.Units[data.Y1][data.X1].GetUnitType()]; !ok {
		return message.GameErrorMessage(ErrMsgUnitCmdNotAllowed), false
	}
	if gl.Units[data.Y1][data.X1].GetUnitStateBit(objects.UnitStateBitMoved) { // has this unit moved?
		return message.GameErrorMessage(ErrMsgUnitAlreadyMoved), false
	}
	if !gl.GridEngine.ValidateMove(data.Y1, data.X1, data.Y2, data.X2) {
		return message.GameErrorMessage(ErrMsgInvalidMove), false
	}
	okAtk, distAtk := gl.GridEngine.ValidateAttack(data.Y2, data.X2, data.YT, data.XT, gl.Units[data.Y1][data.X1])
	if !okAtk {
		return message.GameErrorMessage(ErrMsgInvalidAttack), false
	}
	gl.Units[data.Y2][data.X2], gl.Units[data.Y1][data.X1] = gl.Units[data.Y1][data.X1], gl.Units[data.Y2][data.X2]
	gl.Units[data.Y2][data.X2].ToggleUnitStateBit(objects.UnitStateBitMoved)
	combat.NormalCombat(gl.Units[data.Y2][data.X2], gl.Units[data.YT][data.XT], distAtk)
	replyMsg := &message.GameMessage{
		Cmd:    msg.Cmd,
		Sender: msg.Sender,
		Data: &message.UnitMoveAndAttackMessageDataExt{
			Y1:    data.Y1,
			X1:    data.X1,
			Y2:    data.Y2,
			X2:    data.X2,
			YT:    data.YT,
			XT:    data.XT,
			HPAtk: gl.Units[data.Y2][data.X2].GetUnitHP(),
			HPDef: gl.Units[data.YT][data.XT].GetUnitHP(),
		},
	}
	gl.checkUnitAlive(data.Y2, data.X2)
	gl.checkUnitAlive(data.YT, data.XT)
	return replyMsg, true
}
