package message

// UnitMoveMessageData is message data for CmdUnitMove
type UnitMoveMessageData struct {
	Y1 uint8 `json:"y_1"`
	X1 uint8 `json:"x_1"`
	Y2 uint8 `json:"y_2"`
	X2 uint8 `json:"x_2"`
}
