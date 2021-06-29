package loader

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/cmdwhitelist"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
)

func (gl *GameLoader) handleUnitSwap(msg *message.GameMessage) (*message.GameMessage, bool) {
	data := msg.Data.(*message.UnitSwapMessageData)
	if errMsg := gl.validateUnitOwned(data.Y1, data.X1); len(errMsg) > 0 {
		return message.GameErrorMessage(errMsg), false
	}
	if _, ok := cmdwhitelist.UnitSwapMap[gl.Units[data.Y1][data.X1].UnitType()]; !ok {
		return message.GameErrorMessage(errMsgUnitCmdNotAllowed), false
	}
	if gl.Units[data.Y1][data.X1].GetStateBit(objects.UnitStateBitMoved) { // has this unit moved?
		return message.GameErrorMessage(errMsgUnitAlreadyMoved), false
	}
	if !gl.GridEngine.ValidateSwap(data.Y1, data.X1, data.Y2, data.X2) {
		return message.GameErrorMessage(errMsgInvalidMove), false
	}
	gl.Units[data.Y2][data.X2], gl.Units[data.Y1][data.X1] = gl.Units[data.Y1][data.X1], gl.Units[data.Y2][data.X2]
	gl.Units[data.Y2][data.X2].ToggleStateBit(objects.UnitStateBitMoved)
	return msg, true
}
